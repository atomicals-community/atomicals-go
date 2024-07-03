package atomicals

import (
	"encoding/hex"
	"sort"

	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) transferNft(operation *witness.WitnessAtomicalsOperation, tx *btcjson.TxRawResult) error {
	if operation.IsSplatOperation() { // calculate_nft_atomicals_splat
		atomicalsNfts := make([]*postsql.UTXONftInfo, 0)
		for _, vin := range tx.Vin {
			preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
			if !m.bloomFilter.TestNftLocationID(preNftLocationID) {
				continue
			}
			preNfts, err := m.NftUTXOsByLocationID(preNftLocationID)
			if err != nil {
				log.Log.Panicf("NftUTXOsByLocationID err:%v", err)
			}
			if preNfts == nil {
				continue
			}
			if len(preNfts) == 0 {
				continue
			}
			atomicalsNfts = append(atomicalsNfts, preNfts...)
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
				utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
				utils.IsUnspendableLegacy(scriptPubKeyBytes) {
				output_index = int64(0)
			}
			newUserPk := tx.Vout[output_index].ScriptPubKey.Address
			newLocationID := utils.AtomicalsID(tx.Txid, output_index)
			m.bloomFilter.AddNftLocationID(newLocationID)
			m.UpdateBloomFilter(postsql.FtLocationFilter, m.bloomFilter.Filter[postsql.FtLocationFilter])
			if err := m.TransferNftUTXO(nft.LocationID, newLocationID, newUserPk); err != nil {
				log.Log.Panicf("TransferNftUTXO err:%v", err)
			}
			expected_output_index_incrementing += 1
		}
	} else { // build_nft_input_idx_to_atomical_map && calculate_nft_atomicals_regular
		if utils.IsDmintActivated(operation.RevealLocationHeight) {
			input_idx_to_atomical_ids_map := make(map[int64][]*postsql.UTXONftInfo, 0) // key txInIndex
			for vinIndex, vin := range tx.Vin {
				preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
				if !m.bloomFilter.TestNftLocationID(preNftLocationID) {
					continue
				}
				preNfts, err := m.NftUTXOsByLocationID(preNftLocationID)
				if err != nil {
					log.Log.Panicf("NftUTXOsByLocationID err:%v", err)
				}
				if preNfts == nil {
					continue
				}
				if len(preNfts) == 0 {
					continue
				}
				input_idx_to_atomical_ids_map[int64(vinIndex)] = preNfts
			}
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
					utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
					utils.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					expected_output_index = int64(0)
				}
				for _, nft := range nfts {
					newUserPk := tx.Vout[expected_output_index].ScriptPubKey.Address
					newLocationID := utils.AtomicalsID(tx.Txid, expected_output_index)
					m.bloomFilter.AddNftLocationID(newLocationID)
					m.UpdateBloomFilter(postsql.FtLocationFilter, m.bloomFilter.Filter[postsql.FtLocationFilter])
					if err := m.TransferNftUTXO(nft.LocationID, newLocationID, newUserPk); err != nil {
						log.Log.Panicf("TransferNftUTXO err:%v", err)
					}
				}
				if found_atomical_at_input {
					next_output_idx++
				}
			}
		} else { // calculate_nft_output_index_legacy
			for vinIndex, vin := range tx.Vin {
				preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
				if !m.bloomFilter.TestNftLocationID(preNftLocationID) {
					continue
				}
				preNfts, err := m.NftUTXOsByLocationID(preNftLocationID)
				if err != nil {
					log.Log.Panicf("NftUTXOsByLocationID err:%v", err)
				}
				if preNfts == nil {
					continue
				}
				if len(preNfts) == 0 {
					continue
				}
				expected_output_index := int64(vinIndex)
				// # Assign NFTs the legacy way with 1:1 inputs to outputs
				// # If it was unspendable output, then just set it to the 0th location
				// # ...and never allow an NFT atomical to be burned accidentally by having insufficient number of outputs either
				// # The expected output index will become the 0'th index if the 'x' extract operation was specified or there are insufficient outputs
				// # If this was the 'split' (y) command, then also move them to the 0th output
				if expected_output_index >= int64(len(tx.Vout)) {
					expected_output_index = int64(0)
				}
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[expected_output_index].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
					utils.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					expected_output_index = int64(0)
				}
				for _, nft := range preNfts {
					newUserPk := tx.Vout[expected_output_index].ScriptPubKey.Address
					newLocationID := utils.AtomicalsID(tx.Txid, expected_output_index)
					m.bloomFilter.AddNftLocationID(newLocationID)
					m.UpdateBloomFilter(postsql.FtLocationFilter, m.bloomFilter.Filter[postsql.FtLocationFilter])
					if err := m.TransferNftUTXO(nft.LocationID, newLocationID, newUserPk); err != nil {
						log.Log.Panicf("TransferNftUTXO err:%v", err)
					}
				}
			}
		}
	}
	return nil
}
