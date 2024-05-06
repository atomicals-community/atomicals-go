package atomicals

import (
	"encoding/hex"
	"sort"

	"github.com/atomicals-go/atomicals-core/common"
	"github.com/atomicals-go/atomicals-core/repo/postsql"
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/log"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) transferFtPartialSpliting(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) error {
	if operation.IsSplitOperation() { // color_ft_atomicals_split
		ftAtomicals := make([]*postsql.UTXOFtInfo, 0)
		for _, vin := range tx.Vin {
			preLocationID := common.AtomicalsID(vin.Txid, int64(vin.Vout))
			preFts, err := m.FtUTXOsByLocationID(preLocationID)
			if err != nil {
				log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
			}
			if preFts == nil {
				continue
			}
			if len(preFts) == 0 {
				continue
			}
			ftAtomicals = append(ftAtomicals, preFts...)
			if err := m.DeleteFtUTXO(preLocationID); err != nil {
				log.Log.Panicf("DeleteFtUTXO err:%v", err)
			}
		}
		sort.Slice(ftAtomicals, func(i, j int) bool {
			return ftAtomicals[i].AtomicalsID < ftAtomicals[j].AtomicalsID
		})
		// Todo: haven't catched this param from operation
		for _, ft := range ftAtomicals {
			total_amount_to_skip_potential := operation.Payload.TotalAmountToSkipPotential[ft.LocationID]
			remainingValue := ft.Amount
			for outputIndex, vout := range tx.Vout {
				toBeColoredAmount := int64(vout.Value * common.Satoshi)
				if 0 < total_amount_to_skip_potential {
					total_amount_to_skip_potential -= toBeColoredAmount
					continue
				}
				if remainingValue >= toBeColoredAmount {
					remainingValue -= toBeColoredAmount
				} else {
					// partial colored utxo
					toBeColoredAmount = remainingValue
					remainingValue -= toBeColoredAmount
				}
				locationID := common.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				if err := m.InsertFtUTXO(&postsql.UTXOFtInfo{
					UserPk:          vout.ScriptPubKey.Address,
					AtomicalsID:     ft.AtomicalsID,
					LocationID:      locationID,
					MintTicker:      ft.MintTicker,
					Nonce:           ft.Nonce,
					Time:            ft.Time,
					Bitworkc:        ft.Bitworkc,
					Bitworkr:        ft.Bitworkr,
					MintBitworkVec:  ft.MintBitworkVec,
					MintBitworkcInc: ft.MintBitworkcInc,
					MintBitworkrInc: ft.MintBitworkrInc,
					Amount:          toBeColoredAmount,
				}); err != nil {
					log.Log.Panicf("InsertFtUTXO err:%v", err)
				}
			}
		}
	} else { // color_ft_atomicals_regular
		// a ft in Vin has a total amount: entity.Amount,
		// color exactly the same amount of vout
		// burn the rest ft
		atomicalsFts := make([]*postsql.UTXOFtInfo, 0)
		if common.IsDmintActivated(operation.RevealLocationHeight) {
			atomicalsFtsVinIndexMap := make(map[int64][]*postsql.UTXOFtInfo, 0) // key: vinIndex
			for vinIndex, vin := range tx.Vin {
				preLocationID := common.AtomicalsID(vin.Txid, int64(vin.Vout))
				preFts, err := m.FtUTXOsByLocationID(preLocationID)
				if err != nil {
					log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
				}
				if preFts == nil {
					continue
				}
				if len(preFts) == 0 {
					continue
				}
				atomicalsFtsVinIndexMap[int64(vinIndex)] = preFts
				if err := m.DeleteFtUTXO(preLocationID); err != nil {
					log.Log.Panicf("DeleteFtUTXO err:%v", err)
				}
			}
			seenFtmap := make(map[string]bool, 0) // key: atomicalsID
			for _, fts := range atomicalsFtsVinIndexMap {
				sort.Slice(fts, func(i, j int) bool {
					return fts[i].AtomicalsID < fts[j].AtomicalsID
				})
				for _, ft := range fts {
					if _, ok := seenFtmap[ft.AtomicalsID]; ok {
						continue
					}
					seenFtmap[ft.AtomicalsID] = true
					atomicalsFts = append(atomicalsFts, ft)
				}
			}
		} else {
			for _, vin := range tx.Vin {
				preLocationID := common.AtomicalsID(vin.Txid, int64(vin.Vout))
				preFts, err := m.FtUTXOsByLocationID(preLocationID)
				if err != nil {
					log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
				}
				if preFts == nil {
					continue
				}
				if len(preFts) == 0 {
					continue
				}
				atomicalsFts = append(atomicalsFts, preFts...)
				if err := m.DeleteFtUTXO(preLocationID); err != nil {
					log.Log.Panicf("DeleteFtUTXO err:%v", err)
				}
			}
			sort.Slice(atomicalsFts, func(i, j int) bool {
				return atomicalsFts[i].AtomicalsID < atomicalsFts[j].AtomicalsID
			})
		}

		// calculate_outputs_to_color_for_ft_atomical_ids
		newFts := make([]*postsql.UTXOFtInfo, 0)
		// case 1: cleanly assigned
		nextStartOutIndex := int64(0)
		for _, ft := range atomicalsFts {
			remainingValue := ft.Amount
			for outputIndex := nextStartOutIndex; outputIndex < int64(len(tx.Vout)); outputIndex++ {
				vout := tx.Vout[outputIndex]
				if int64(outputIndex) < nextStartOutIndex {
					continue
				}
				toBeColoredAmount := int64(vout.Value * common.Satoshi)
				if remainingValue >= toBeColoredAmount {
					remainingValue -= toBeColoredAmount
				} else {
					// partial colored utxo
					toBeColoredAmount = remainingValue
					remainingValue -= toBeColoredAmount
				}
				nextStartOutIndex = outputIndex + 1
				scriptPubKeyBytes, err := hex.DecodeString(vout.ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if common.IsUnspendableGenesis(scriptPubKeyBytes) ||
					common.IsUnspendableLegacy(scriptPubKeyBytes) {
					continue
				}
				newLocationID := common.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				newFts = append(newFts, &postsql.UTXOFtInfo{
					UserPk:          vout.ScriptPubKey.Address,
					AtomicalsID:     ft.AtomicalsID,
					LocationID:      newLocationID,
					MintTicker:      ft.MintTicker,
					Nonce:           ft.Nonce,
					Time:            ft.Time,
					Bitworkc:        ft.Bitworkc,
					Bitworkr:        ft.Bitworkr,
					MintBitworkVec:  ft.MintBitworkVec,
					MintBitworkcInc: ft.MintBitworkcInc,
					MintBitworkrInc: ft.MintBitworkrInc,
					Amount:          int64(vout.Value * common.Satoshi),
				})
			}
			// # If the output slots did not fit cleanly, then default to just assigning everything from the 0'th output index
			if remainingValue == 0 {
				for _, ft := range newFts {
					if err := m.InsertFtUTXO(ft); err != nil {
						log.Log.Panicf("InsertFtUTXO err:%v", err)
					}
				}
				return nil
			}
		}
		// case 2: not cleanly assigned
		newFts = make([]*postsql.UTXOFtInfo, 0)
		for _, ft := range atomicalsFts {
			remainingValue := ft.Amount
			for outputIndex, vout := range tx.Vout {
				toBeColoredAmount := int64(vout.Value * common.Satoshi)
				if remainingValue >= toBeColoredAmount {
					remainingValue -= toBeColoredAmount
				} else {
					// partial colored utxo
					toBeColoredAmount = remainingValue
					remainingValue -= toBeColoredAmount
				}
				scriptPubKeyBytes, err := hex.DecodeString(vout.ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if common.IsUnspendableGenesis(scriptPubKeyBytes) ||
					common.IsUnspendableLegacy(scriptPubKeyBytes) {
					continue
				}
				newLocationID := common.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				newFts = append(newFts, &postsql.UTXOFtInfo{
					UserPk:          vout.ScriptPubKey.Address,
					AtomicalsID:     ft.AtomicalsID,
					LocationID:      newLocationID,
					MintTicker:      ft.MintTicker,
					Nonce:           ft.Nonce,
					Time:            ft.Time,
					Bitworkc:        ft.Bitworkc,
					Bitworkr:        ft.Bitworkr,
					MintBitworkVec:  ft.MintBitworkVec,
					MintBitworkcInc: ft.MintBitworkcInc,
					MintBitworkrInc: ft.MintBitworkrInc,
					Amount:          int64(vout.Value * common.Satoshi),
				})
			}
			// # If the output slots did not fit cleanly, then default to just assigning everything from the 0'th output index
			for _, ft := range newFts {
				if err := m.InsertFtUTXO(ft); err != nil {
					log.Log.Panicf("InsertFtUTXO err:%v", err)
				}
			}
			return nil
		}

	}
	return nil
}
