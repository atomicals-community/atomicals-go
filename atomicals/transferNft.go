package atomicals

import (
	"encoding/hex"

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
		expected_output_index_incrementing := int64(0) // # Begin assigning splatted atomicals at the 0'th index
		for _, vin := range tx.Vin {
			preAtomicalsID := atomicalsID(vin.Txid, int64(vin.Vout))
			if _, ok := m.UTXOs[preAtomicalsID]; !ok {
				continue
			}
			preAtomicals := m.UTXOs[preAtomicalsID].Nft
			if preAtomicals != nil {
				continue
			}
			for _, entity := range preAtomicals {
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
				entity.UserPk = tx.Vout[output_index].ScriptPubKey.Address
				atomicalsID := atomicalsID(operation.RevealLocationTxID, output_index)
				m.UTXOs[atomicalsID].AtomicalID = atomicalsID
				m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, entity)
				expected_output_index_incrementing += 1
			}
			m.UTXOs[preAtomicalsID].Nft = nil
		}
	} else { // calculate_nft_atomicals_regular
		if common.IsDmintActivated(operation.RevealLocationHeight) {
			expected_output_index_incrementing := int64(0)
			for _, vin := range tx.Vin {
				preAtomicalsID := atomicalsID(vin.Txid, int64(vin.Vout))
				if _, ok := m.UTXOs[preAtomicalsID]; !ok {
					continue
				}
				preAtomicals := m.UTXOs[preAtomicalsID].Nft
				if preAtomicals != nil {
					continue
				}
				output_index := expected_output_index_incrementing
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[output_index].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if output_index >= int64(len(tx.Vout)) ||
					common.IsUnspendableGenesis(scriptPubKeyBytes) ||
					common.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					output_index = int64(0)
				}
				for _, entity := range preAtomicals {
					entity.UserPk = tx.Vout[output_index].ScriptPubKey.Address
				}
				atomicalsID := atomicalsID(operation.RevealLocationTxID, output_index)
				m.UTXOs[atomicalsID].AtomicalID = atomicalsID
				m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, preAtomicals...)
				expected_output_index_incrementing += 1
				m.UTXOs[preAtomicalsID].Nft = nil
			}
		} else { // calculate_nft_output_index_legacy
			for vinIndex, vin := range tx.Vin {
				preAtomicalsID := atomicalsID(vin.Txid, int64(vin.Vout))
				if _, ok := m.UTXOs[preAtomicalsID]; !ok {
					continue
				}
				preAtomicals := m.UTXOs[preAtomicalsID].Nft
				if preAtomicals != nil {
					continue
				}
				output_index := int64(vinIndex)
				//   # Assign NFTs the legacy way with 1:1 inputs to outputs
				// # If it was unspendable output, then just set it to the 0th location
				// # ...and never allow an NFT atomical to be burned accidentally by having insufficient number of outputs either
				// # The expected output index will become the 0'th index if the 'x' extract operation was specified or there are insufficient outputs
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[output_index].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				// # If this was the 'split' (y) command, then also move them to the 0th output
				if output_index >= int64(len(tx.Vout)) ||
					common.IsUnspendableGenesis(scriptPubKeyBytes) ||
					common.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					output_index = int64(0)
				}
				for _, entity := range preAtomicals {
					entity.UserPk = tx.Vout[output_index].ScriptPubKey.Address
				}
				atomicalsID := atomicalsID(operation.RevealLocationTxID, output_index)
				m.UTXOs[atomicalsID].AtomicalID = atomicalsID
				m.UTXOs[atomicalsID].Nft = append(m.UTXOs[atomicalsID].Nft, preAtomicals...)
				m.UTXOs[preAtomicalsID].Nft = nil
			}
		}
	}
	return nil
}
