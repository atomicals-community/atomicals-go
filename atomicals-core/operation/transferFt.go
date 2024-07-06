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

func (m *Atomicals) transferFt(operation *witness.WitnessAtomicalsOperation, tx *btcjson.TxRawResult) error {
	if operation.IsSplitOperation() { // color_ft_atomicals_split
		ftAtomicals := make([]*postsql.UTXOFtInfo, 0)
		for _, vin := range tx.Vin {
			preLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
			preFts, err := m.FtUTXOsByLocationID(preLocationID)
			if err != nil {
				log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
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
		for _, ft := range ftAtomicals {
			totalAmountToSkipPotential := operation.Payload.Args.TotalAmountToSkipPotential[ft.LocationID]
			remainingValue := ft.Amount
			for outputIndex, vout := range tx.Vout {
				if 0 < totalAmountToSkipPotential {
					totalAmountToSkipPotential -= int64(vout.Value * utils.Satoshi)
					continue
				}
				if remainingValue < int64(vout.Value*utils.Satoshi) { // burn rest ft
					break
				}
				remainingValue -= int64(vout.Value * utils.Satoshi)
				locationID := utils.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				panic("locationID")
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
					Amount:          int64(vout.Value * utils.Satoshi),
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
		if utils.IsDmintActivated(operation.RevealLocationHeight) {
			atomicalsFtsVinIndexMap := make(map[int64][]*postsql.UTXOFtInfo, 0) // key: vinIndex
			for vinIndex, vin := range tx.Vin {
				preLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
				preFts, err := m.FtUTXOsByLocationID(preLocationID)
				if err != nil {
					log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
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
				preLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
				preFts, err := m.FtUTXOsByLocationID(preLocationID)
				if err != nil {
					log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
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
		nextStartOutIdx := int64(0)
		nonCleanOutputSlots := false
		for _, ft := range atomicalsFts {
			cleanly_assigned, assignedVoutIndex, fts := assignExpectedOutputsBasic(ft, operation, tx, nextStartOutIdx)
			if cleanly_assigned {
				nextStartOutIdx = assignedVoutIndex + 1
				newFts = append(newFts, fts...)
			} else {
				nonCleanOutputSlots = true
				newFts = make([]*postsql.UTXOFtInfo, 0)
				break
			}
		}
		// # If the output slots did not fit cleanly, then default to just assigning everything from the 0'th output index
		if nonCleanOutputSlots {
			newFts = make([]*postsql.UTXOFtInfo, 0)
			for _, ft := range atomicalsFts {
				_, _, fts := assignExpectedOutputsBasic(ft, operation, tx, 0) //always is 0'th
				newFts = append(newFts, fts...)
			}
		}
		for _, ft := range newFts {
			if err := m.InsertFtUTXO(ft); err != nil {
				log.Log.Panicf("InsertFtUTXO err:%v", err)
			}
		}
		return nil
	}
	return nil
}

// assign_expected_outputs_basic
func assignExpectedOutputsBasic(ft *postsql.UTXOFtInfo, operation *witness.WitnessAtomicalsOperation, tx *btcjson.TxRawResult, startOutIdx int64) (bool, int64, []*postsql.UTXOFtInfo) {
	newFts := make([]*postsql.UTXOFtInfo, 0)
	remainingValue := ft.Amount
	if startOutIdx >= int64(len(tx.Vout)) {
		return false, startOutIdx, nil
	}
	assignedVoutIndex := int64(0)
	for outputIndex, vout := range tx.Vout {
		if int64(outputIndex) < startOutIdx {
			continue
		}
		assignedVoutIndex = int64(outputIndex)
		scriptPubKeyBytes, err := hex.DecodeString(vout.ScriptPubKey.Hex)
		if err != nil {
			panic(err)
		}
		if utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
			utils.IsUnspendableLegacy(scriptPubKeyBytes) {
			continue
		}
		if int64(vout.Value*utils.Satoshi) > remainingValue { // burn rest ft
			return false, assignedVoutIndex, nil
		}
		remainingValue -= int64(vout.Value * utils.Satoshi)
		locationID := utils.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
		newFts = append(newFts, &postsql.UTXOFtInfo{
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
			Amount:          int64(vout.Value * utils.Satoshi),
		})
		if remainingValue == 0 {
			return true, assignedVoutIndex, newFts
		}
	}
	return false, assignedVoutIndex, nil
}
