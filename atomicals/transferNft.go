package atomicals

import (
	"encoding/hex"
	"sort"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) transferNft(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) error {
	if operation.IsSplatOperation() { // calculate_nft_atomicals_splat
		// # Splat takes all of the NFT atomicals across all inputs (including multiple atomicals at the same utxo)
		// # and then separates them into their own distinctive output such that the result of the operation is no two atomicals
		// # will share a resulting output. This operation requires that there are at least as many outputs as there are NFT atomicals
		// # If there are not enough, then this is considered a noop and those extra NFTs are assigned to output 0
		// # If there are enough outputs, then the earliest atomical (sorted lexicographically in ascending order) goes to the 0'th output,
		// # then the second atomical goes to the 1'st output, etc until all atomicals are assigned to their own output.
		atomicalsNfts := make([]*UserNftInfo, 0)
		for _, vin := range tx.Vin {
			preNftLocationID := atomicalsID(vin.Txid, int64(vin.Vout))
			if _, ok := m.NftUTXOsByLocationID[preNftLocationID]; !ok {
				continue
			}
			preNfts := m.NftUTXOsByLocationID[preNftLocationID]
			if preNfts != nil || len(preNfts) == 0 {
				continue
			}
			atomicalsNfts = append(atomicalsNfts, preNfts...)
			m.NftUTXOsByLocationID[preNftLocationID] = nil
		}
		sort.Slice(atomicalsNfts, func(i, j int) bool {
			return atomicalsNfts[i].AtomicalsID < atomicalsNfts[j].AtomicalsID
		})
		expected_output_index_incrementing := int64(0) // # Begin assigning splatted atomicals at the 0'th index
		for _, nft := range atomicalsNfts {
			output_index := expected_output_index_incrementing
			scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[output_index].ScriptPubKey.Hex)
			if err != nil {
				panic("err")
			}
			if output_index >= int64(len(tx.Vout)) ||
				common.IsUnspendableGenesis(scriptPubKeyBytes) ||
				common.IsUnspendableLegacy(scriptPubKeyBytes) {
				output_index = int64(0)
			}
			nft.UserPk = tx.Vout[output_index].ScriptPubKey.Address
			locationID := atomicalsID(operation.RevealLocationTxID, output_index)
			nft.LocationID = locationID
			m.NftUTXOsByLocationID[locationID] = append(m.NftUTXOsByLocationID[locationID], nft)
			expected_output_index_incrementing += 1
		}
	} else { // build_nft_input_idx_to_atomical_map && calculate_nft_atomicals_regular
		input_idx_to_atomical_ids_map := make(map[int64][]*UserNftInfo, 0) // key txInIndex
		for vinIndex, vin := range tx.Vin {
			preNftLocationID := atomicalsID(vin.Txid, int64(vin.Vout))
			if _, ok := m.NftUTXOsByLocationID[preNftLocationID]; !ok {
				continue
			}
			preNfts := m.NftUTXOsByLocationID[preNftLocationID]
			if preNfts != nil || len(preNfts) == 0 {
				continue
			}
			input_idx_to_atomical_ids_map[int64(vinIndex)] = preNfts
			m.NftUTXOsByLocationID[preNftLocationID] = nil
		}
		if common.IsDmintActivated(operation.RevealLocationHeight) {
			next_output_idx := int64(0)
			found_atomical_at_input := false
			for _, nfts := range input_idx_to_atomical_ids_map {
				found_atomical_at_input = true
				expected_output_index := next_output_idx
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[expected_output_index].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if expected_output_index >= int64(len(tx.Vout)) ||
					common.IsUnspendableGenesis(scriptPubKeyBytes) ||
					common.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					expected_output_index = int64(0)
				}
				for _, nft := range nfts {
					nft.UserPk = tx.Vout[expected_output_index].ScriptPubKey.Address
					locationID := atomicalsID(operation.RevealLocationTxID, expected_output_index)
					nft.LocationID = locationID
					m.NftUTXOsByLocationID[locationID] = append(m.NftUTXOsByLocationID[locationID], nft)
				}
			}
			if found_atomical_at_input {
				next_output_idx++
			}
		} else { // calculate_nft_output_index_legacy
			for vinIndex, vin := range tx.Vin {
				preAtomicalsID := atomicalsID(vin.Txid, int64(vin.Vout))
				if _, ok := m.NftUTXOsByLocationID[preAtomicalsID]; !ok {
					continue
				}
				preAtomicalsNft := m.NftUTXOsByLocationID[preAtomicalsID]
				if preAtomicalsNft != nil || len(preAtomicalsNft) == 0 {
					continue
				}
				expected_output_index := int64(vinIndex)
				//   # Assign NFTs the legacy way with 1:1 inputs to outputs
				// # If it was unspendable output, then just set it to the 0th location
				// # ...and never allow an NFT atomical to be burned accidentally by having insufficient number of outputs either
				// # The expected output index will become the 0'th index if the 'x' extract operation was specified or there are insufficient outputs
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[expected_output_index].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				// # If this was the 'split' (y) command, then also move them to the 0th output
				if expected_output_index >= int64(len(tx.Vout)) ||
					common.IsUnspendableGenesis(scriptPubKeyBytes) ||
					common.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					expected_output_index = int64(0)
				}
				for _, nft := range preAtomicalsNft {
					nft.UserPk = tx.Vout[expected_output_index].ScriptPubKey.Address
					locationID := atomicalsID(operation.RevealLocationTxID, expected_output_index)
					nft.LocationID = locationID
					m.NftUTXOsByLocationID[locationID] = append(m.NftUTXOsByLocationID[locationID], nft)
				}
				m.NftUTXOsByLocationID[preAtomicalsID] = nil
			}
		}
	}
	return nil
}
