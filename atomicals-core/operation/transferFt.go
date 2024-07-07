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
					Nonce:           ft.Nonce,
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
					deleteFts = append(deleteFts, ft)
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
				deleteFts = append(deleteFts, preFts...)
			}
			sort.Slice(deleteFts, func(i, j int) bool {
				return deleteFts[i].AtomicalsID < deleteFts[j].AtomicalsID
			})
		}

		// calculate_outputs_to_color_for_ft_atomical_ids
		voutRemainingSpace := make([]int64, len(tx.Vout))
		for i, vout := range tx.Vout {
			voutRemainingSpace[i] = utils.MulSatoshi(vout.Value)
		}
		newFtAmount := int64(0)
		outputIndex := int64(0)
		for i, ft := range deleteFts {
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
					continue
				}
				locationID := utils.AtomicalsID(operation.RevealLocationTxID, int64(outputIndex))
				if voutRemainingSpace[outputIndex] > ft.Amount { // burn rest ft
					log.Log.Infof("height:%v %v, time:%v", i, voutRemainingSpace[outputIndex], ft.Amount)
					voutRemainingSpace[outputIndex] = voutRemainingSpace[outputIndex] - ft.Amount
					newFtAmount += ft.Amount
					break
				} else if voutRemainingSpace[outputIndex] == ft.Amount { // burn rest ft
					voutRemainingSpace[outputIndex] = voutRemainingSpace[outputIndex] - ft.Amount
					newFtAmount += ft.Amount
					outputIndex = outputIndex + 1
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
						Amount:          newFtAmount,
					})
					newFtAmount = 0
					break
				} else if voutRemainingSpace[outputIndex] < ft.Amount {
					voutRemainingSpace[outputIndex] = 0
					outputIndex = outputIndex + 1
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
						Amount:          utils.MulSatoshi(vout.Value),
					})
					ft.Amount -= utils.MulSatoshi(vout.Value)
					newFtAmount = 0
				}
			}
		}
	}
	if tx.Txid == "ec4d4c0b24196f838e059927c84b42b87d91dc63a9692392bdfe8c78b89702e9" {
		log.Log.Infof("Txid:%v", tx.Txid)
	}
	if len(newFts) != 0 {
		for _, v := range newFts {
			log.Log.Infof("Amount:%v", v.Amount)
		}
		log.Log.Infof("Txid:%v", tx.Txid)
		log.Log.Infof("Txid:%v", tx.Txid)
	}
	return deleteFts, newFts, nil
}
