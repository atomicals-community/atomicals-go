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

func (m *Atomicals) transferFt(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) (deleteFts []*postsql.UTXOFtInfo, newFts []*postsql.UTXOFtInfo, err error) {
	if operation.IsSplitOperation() { // color_ft_atomicals_split
		for _, vin := range tx.Vin {
			preLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
			preFts, err := m.FtUTXOsByLocationID(preLocationID)
			if err != nil {
				log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
			}
			if len(preFts) == 0 {
				continue
			}
			deleteFts = append(deleteFts, preFts...)
		}
		sort.Slice(deleteFts, func(i, j int) bool {
			return deleteFts[i].AtomicalsID < deleteFts[j].AtomicalsID
		})
		for _, ft := range deleteFts {
			totalAmountToSkipPotential := operation.Payload.Args.TotalAmountToSkipPotential[ft.LocationID]
			remainingValue := ft.Amount
			for outputIndex, vout := range tx.Vout {
				amount := utils.MulSatoshi(vout.Value)
				if 0 < totalAmountToSkipPotential {
					totalAmountToSkipPotential -= amount
					continue
				}
				if remainingValue < amount { // burn rest ft
					break
				}
				remainingValue -= amount
				locationID := utils.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				newFts = append(newFts, &postsql.UTXOFtInfo{
					UserPk:          vout.ScriptPubKey.Address,
					AtomicalsID:     ft.AtomicalsID,
					LocationID:      locationID,
					MintTicker:      ft.MintTicker,
					Time:            ft.Time,
					Bitworkc:        ft.Bitworkc,
					Bitworkr:        ft.Bitworkr,
					MintBitworkVec:  ft.MintBitworkVec,
					MintBitworkcInc: ft.MintBitworkcInc,
					MintBitworkrInc: ft.MintBitworkrInc,
					Amount:          amount,
				})
			}
		}
	} else { // color_ft_atomicals_regular
		// a ft in Vin has a total amount: entity.Amount,
		// color exactly the same amount of vout
		// burn the rest ft
		for _, vin := range tx.Vin {
			preLocationID := utils.AtomicalsID(vin.Txid, int64(vin.Vout))
			preFts, err := m.FtUTXOsByLocationID(preLocationID)
			if err != nil {
				log.Log.Panicf("FtUTXOsByLocationID err:%v", err)
			}
			if len(preFts) == 0 {
				continue
			}
			deleteFts = append(deleteFts, preFts...)
		}
		sort.Slice(deleteFts, func(i, j int) bool {
			return deleteFts[i].AtomicalsID < deleteFts[j].AtomicalsID
		})
		deleteFtMap := make(map[string][]*postsql.UTXOFtInfo, 0)
		for _, ft := range deleteFts {
			if _, ok := deleteFtMap[ft.AtomicalsID]; !ok {
				deleteFtMap[ft.AtomicalsID] = make([]*postsql.UTXOFtInfo, 0)
			}
			deleteFtMap[ft.AtomicalsID] = append(deleteFtMap[ft.AtomicalsID], ft)
		}
		// calculate_outputs_to_color_for_ft_atomical_ids
		newFts = m.AssignFt(tx, operation, deleteFtMap)
	}
	return deleteFts, newFts, nil
}

func (m *Atomicals) AssignFt(tx btcjson.TxRawResult, operation *witness.WitnessAtomicalsOperation, deleteFtMap map[string][]*postsql.UTXOFtInfo) []*postsql.UTXOFtInfo {
	newFts := make([]*postsql.UTXOFtInfo, 0)
	for _, ftSlice := range deleteFtMap {
		voutRemainingSpace := make([]int64, len(tx.Vout))
		for i, vout := range tx.Vout {
			voutRemainingSpace[i] = utils.MulSatoshi(vout.Value)
		}
		newFtAmount := int64(0)
		outputIndex := int64(0)
		for i, ft := range ftSlice {
			ftAmount := ft.Amount
			for {
				if outputIndex >= int64(len(tx.Vout)) {
					break
				}
				vout := tx.Vout[outputIndex]
				scriptPubKeyBytes, err := hex.DecodeString(vout.ScriptPubKey.Hex)
				if err != nil {
					panic(err)
				}
				if utils.IsUnspendableGenesis(scriptPubKeyBytes) ||
					utils.IsUnspendableLegacy(scriptPubKeyBytes) {
					outputIndex = outputIndex + 1
					continue
				}
				locationID := utils.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				if voutRemainingSpace[outputIndex] > ftAmount { // burn rest ft
					voutRemainingSpace[outputIndex] = voutRemainingSpace[outputIndex] - ftAmount
					newFtAmount += ftAmount
					if utils.IsCustomColoring(operation.RevealLocationHeight) && i == (len(ftSlice)-1) {
						newFts = append(newFts, &postsql.UTXOFtInfo{
							UserPk:          vout.ScriptPubKey.Address,
							AtomicalsID:     ft.AtomicalsID,
							LocationID:      locationID,
							MintTicker:      ft.MintTicker,
							Time:            ft.Time,
							Bitworkc:        ft.Bitworkc,
							Bitworkr:        ft.Bitworkr,
							MintBitworkVec:  ft.MintBitworkVec,
							MintBitworkcInc: ft.MintBitworkcInc,
							MintBitworkrInc: ft.MintBitworkrInc,
							Amount:          newFtAmount,
						})
					}
					break
				} else if voutRemainingSpace[outputIndex] == ftAmount { // burn rest ft
					voutRemainingSpace[outputIndex] = voutRemainingSpace[outputIndex] - ftAmount
					newFtAmount += ftAmount
					outputIndex = outputIndex + 1
					newFts = append(newFts, &postsql.UTXOFtInfo{
						UserPk:          vout.ScriptPubKey.Address,
						AtomicalsID:     ft.AtomicalsID,
						LocationID:      locationID,
						MintTicker:      ft.MintTicker,
						Time:            ft.Time,
						Bitworkc:        ft.Bitworkc,
						Bitworkr:        ft.Bitworkr,
						MintBitworkVec:  ft.MintBitworkVec,
						MintBitworkcInc: ft.MintBitworkcInc,
						MintBitworkrInc: ft.MintBitworkrInc,
						Amount:          newFtAmount,
					})
					newFtAmount = 0
					break
				} else if voutRemainingSpace[outputIndex] < ftAmount {
					voutRemainingSpace[outputIndex] = 0
					newFts = append(newFts, &postsql.UTXOFtInfo{
						UserPk:          vout.ScriptPubKey.Address,
						AtomicalsID:     ft.AtomicalsID,
						LocationID:      locationID,
						MintTicker:      ft.MintTicker,
						Time:            ft.Time,
						Bitworkc:        ft.Bitworkc,
						Bitworkr:        ft.Bitworkr,
						MintBitworkVec:  ft.MintBitworkVec,
						MintBitworkcInc: ft.MintBitworkcInc,
						MintBitworkrInc: ft.MintBitworkrInc,
						Amount:          utils.MulSatoshi(vout.Value),
					})
					if outputIndex != int64(len(tx.Vout))-1 {
						ftAmount -= utils.MulSatoshi(vout.Value)
						newFtAmount = 0
					}
					outputIndex = outputIndex + 1
				}
			}
		}
	}
	return newFts
}
