package atomicals

import (
	"encoding/hex"
	"sort"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) transferFt(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) error {
	if operation.IsSplitOperation() { // color_ft_atomicals_split
		// a dmt in Vin has a total amount: entity.Amount,
		// retain total_amount_to_skip_potential number of Dmt from vin
		// color exactly the same amount of vout
		// burn the rest Dmt
		// total_amount_to_skip_potential := float64(operation.Payload.total_amount_to_skip_potential(preAtomicalsID))
		total_amount_to_skip_potential := int64(0) // Todo: haven't catched this param, it confused me, why we need to skip of some amount?
		ftAtomicals := make([]*UserFtInfo, 0)
		for _, vin := range tx.Vin {
			preLocationID := atomicalsID(vin.Txid, int64(vin.Vout))
			if _, ok := m.FtUTXOs[preLocationID]; !ok {
				continue
			}
			preAtomicals := m.FtUTXOs[preLocationID]
			if preAtomicals != nil || len(preAtomicals) == 0 {
				continue
			}
			if 0 < total_amount_to_skip_potential {
				for _, ft := range preAtomicals {
					ftAtomicals = append(ftAtomicals, ft)
					ft.Amount = total_amount_to_skip_potential
				}
			} else {
				ftAtomicals = append(ftAtomicals, preAtomicals...)
				m.FtUTXOs[preLocationID] = nil
			}
		}
		sort.Slice(ftAtomicals, func(i, j int) bool {
			return ftAtomicals[i].AtomicalsID < ftAtomicals[j].AtomicalsID
		})
		for _, ft := range ftAtomicals {
			remaining_value := ft.Amount
			for outputIndex, vout := range tx.Vout {
				if remaining_value < int64(vout.Value*common.Satoshi) { // burn rest ft
					break
				}
				remaining_value -= int64(vout.Value * common.Satoshi)
				locationID := atomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				m.ensureFtUTXONotNil(locationID)
				m.FtUTXOs[locationID] = append(m.FtUTXOs[locationID], &UserFtInfo{
					UserPk:          vout.ScriptPubKey.Address,
					AtomicalsID:     ft.AtomicalsID,
					LocaiontID:      locationID,
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
		}
	} else { // color_ft_atomicals_regular
		// a ft in Vin has a total amount: entity.Amount,
		// color exactly the same amount of vout
		// burn the rest ft
		atomicalsFts := make([]*UserFtInfo, 0)
		if common.IsDmintActivated(operation.RevealLocationHeight) {
			atomicalsFtsVinIndexMap := make(map[int64][]*UserFtInfo, 0) // key: vinIndex
			for vinIndex, vin := range tx.Vin {
				preNftLocationID := atomicalsID(vin.Txid, int64(vin.Vout))
				if _, ok := m.FtUTXOs[preNftLocationID]; !ok {
					continue
				}
				preFts := m.FtUTXOs[preNftLocationID]
				if preFts != nil || len(preFts) == 0 {
					continue
				}
				atomicalsFtsVinIndexMap[int64(vinIndex)] = preFts
				m.FtUTXOs[preNftLocationID] = nil
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
				preNftLocationID := atomicalsID(vin.Txid, int64(vin.Vout))
				if _, ok := m.FtUTXOs[preNftLocationID]; !ok {
					continue
				}
				preFts := m.FtUTXOs[preNftLocationID]
				if preFts != nil || len(preFts) == 0 {
					continue
				}
				atomicalsFts = append(atomicalsFts, preFts...)
				m.FtUTXOs[preNftLocationID] = nil
			}
			sort.Slice(atomicalsFts, func(i, j int) bool {
				return atomicalsFts[i].AtomicalsID < atomicalsFts[j].AtomicalsID
			})
		}

		// calculate_outputs_to_color_for_ft_atomical_ids
		newFts := make([]*UserFtInfo, 0)
		next_start_out_idx := int64(0)
		non_clean_output_slots := false
		for _, ft := range atomicalsFts {
			// assign_expected_outputs_basic
			cleanly_assigned, assignedVoutIndex, fts := assign_expected_outputs_basic(ft, operation, tx, next_start_out_idx)
			if cleanly_assigned {
				next_start_out_idx = assignedVoutIndex + 1
				newFts = append(newFts, fts...)
			} else {
				non_clean_output_slots = true
				newFts = make([]*UserFtInfo, 0)
				break
			}
		}
		// # If the output slots did not fit cleanly, then default to just assigning everything from the 0'th output index
		if non_clean_output_slots {
			for _, ft := range atomicalsFts {
				_, _, fts := assign_expected_outputs_basic(ft, operation, tx, 0) //always is 0'th
				newFts = append(newFts, fts...)
			}
		}
		for _, ft := range newFts {
			m.ensureFtUTXONotNil(ft.LocaiontID)
			m.FtUTXOs[ft.LocaiontID] = append(m.FtUTXOs[ft.LocaiontID], ft)
		}
		return nil
	}
	return nil
}

func assign_expected_outputs_basic(ft *UserFtInfo, operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult, start_out_idx int64) (bool, int64, []*UserFtInfo) {
	newFts := make([]*UserFtInfo, 0)
	remaining_value := ft.Amount
	if start_out_idx >= int64(len(tx.Vout)) {
		return false, start_out_idx, nil
	}
	assignedVoutIndex := int64(0)
	for outputIndex, vout := range tx.Vout {
		if int64(outputIndex) < start_out_idx {
			continue
		}
		assignedVoutIndex = int64(outputIndex)
		scriptPubKeyBytes, err := hex.DecodeString(vout.ScriptPubKey.Hex)
		if err != nil {
			panic(err)
		}
		if common.IsUnspendableGenesis(scriptPubKeyBytes) ||
			common.IsUnspendableLegacy(scriptPubKeyBytes) {
			continue
		}
		if int64(vout.Value*common.Satoshi) > remaining_value { // burn rest ft
			return false, assignedVoutIndex, nil
		}
		remaining_value -= int64(vout.Value * common.Satoshi)
		locationID := atomicalsID(operation.RevealLocationTxID, int64(outputIndex))
		newFts = append(newFts, &UserFtInfo{
			UserPk:          vout.ScriptPubKey.Address,
			AtomicalsID:     ft.AtomicalsID,
			LocaiontID:      locationID,
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
		if remaining_value == 0 {
			return true, assignedVoutIndex, newFts
		}
	}
	return false, assignedVoutIndex, nil
}
