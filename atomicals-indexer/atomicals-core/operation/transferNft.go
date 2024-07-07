package atomicals

import (
	"encoding/hex"
	"sort"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) transferNft(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) (updateNfts []*postsql.UTXONftInfo, err error) {
	if operation.IsSplatOperation() { // calculate_nft_atomicals_splat
		atomicalsNfts := make([]*postsql.UTXONftInfo, 0)
		for _, vin := range tx.Vin {
			preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
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
		expectedOutputIndexIncrementing := int64(0) // # Begin assigning splatted atomicals at the 0'th index
		for _, nft := range atomicalsNfts {
			outputIndex := expectedOutputIndexIncrementing
			scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[outputIndex].ScriptPubKey.Hex)
			if err != nil {
				panic("err")
			}
			if outputIndex >= int64(len(tx.Vout)) ||
				utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
				utils.IsUnspendableLegacy(scriptPubKeyBytes) {
				outputIndex = int64(0)
			}
			nft.LocationID = utils.AtomicalsID(tx.Txid, outputIndex)
			nft.UserPk = tx.Vout[outputIndex].ScriptPubKey.Address
			updateNfts = append(updateNfts, nft)
			expectedOutputIndexIncrementing += 1
		}
	} else { // build_nft_input_idx_to_atomical_map && calculate_nft_atomicals_regular
		if utils.IsDmintActivated(operation.RevealLocationHeight) {
			inputIdxToAtomicalIdsMap := make(map[int64][]*postsql.UTXONftInfo, 0) // key txInIndex
			for vinIndex, vin := range tx.Vin {
				preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
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
				inputIdxToAtomicalIdsMap[int64(vinIndex)] = preNfts
			}
			nextOutputIdx := int64(0)
			foundAtomicalAtInput := false
			for _, nfts := range inputIdxToAtomicalIdsMap {
				foundAtomicalAtInput = true
				expectedOutputIndex := nextOutputIdx
				if expectedOutputIndex >= int64(len(tx.Vout)) {
					expectedOutputIndex = int64(0)
				}
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[expectedOutputIndex].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
					utils.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					expectedOutputIndex = int64(0)
				}
				for _, nft := range nfts {
					nft.LocationID = utils.AtomicalsID(tx.Txid, expectedOutputIndex)
					nft.UserPk = tx.Vout[expectedOutputIndex].ScriptPubKey.Address
					updateNfts = append(updateNfts, nft)
				}
				if foundAtomicalAtInput {
					nextOutputIdx++
				}
			}
		} else { // calculate_nft_output_index_legacy
			for vinIndex, vin := range tx.Vin {
				preNftLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
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
				expectedOutputIndex := int64(vinIndex)
				// # Assign NFTs the legacy way with 1:1 inputs to outputs
				// # If it was unspendable output, then just set it to the 0th location
				// # ...and never allow an NFT atomical to be burned accidentally by having insufficient number of outputs either
				// # The expected output index will become the 0'th index if the 'x' extract operation was specified or there are insufficient outputs
				// # If this was the 'split' (y) command, then also move them to the 0th output
				if expectedOutputIndex >= int64(len(tx.Vout)) {
					expectedOutputIndex = int64(0)
				}
				scriptPubKeyBytes, err := hex.DecodeString(tx.Vout[expectedOutputIndex].ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
					utils.IsUnspendableLegacy(scriptPubKeyBytes) ||
					operation.IsSplitOperation() {
					expectedOutputIndex = int64(0)
				}
				for _, nft := range preNfts {
					nft.LocationID = utils.AtomicalsID(tx.Txid, expectedOutputIndex)
					nft.UserPk = tx.Vout[expectedOutputIndex].ScriptPubKey.Address
					updateNfts = append(updateNfts, nft)
				}
			}
		}
	}
	return updateNfts, nil
}
