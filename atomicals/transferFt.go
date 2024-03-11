package atomicals

import (
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
		for _, vin := range tx.Vin {
			preAtomicalsID := atomicalsID(vin.Txid, int64(vin.Vout))
			if _, ok := m.UTXOs[preAtomicalsID]; !ok {
				continue
			}
			// total_amount_to_skip_potential := float64(operation.Payload.total_amount_to_skip_potential(preAtomicalsID))
			total_amount_to_skip_potential := float64(0) // Todo: haven't catched this param
			for index, entity := range m.UTXOs[preAtomicalsID].DistributedFt {
				remaining_value := entity.Amount
				if 0 < total_amount_to_skip_potential {
					m.UTXOs[preAtomicalsID].DistributedFt[index].Amount = total_amount_to_skip_potential
				} else {
					m.UTXOs[preAtomicalsID].DistributedFt = append(m.UTXOs[preAtomicalsID].DistributedFt[:index], m.UTXOs[preAtomicalsID].DistributedFt[index+1:]...)
				}
				for output_index, vout := range tx.Vout {
					if 0 < total_amount_to_skip_potential {
						total_amount_to_skip_potential -= vout.Value * common.Satoshi
						continue
					}
					if remaining_value < vout.Value*common.Satoshi { // if so, this Dmt will be burned
						break
					}
					remaining_value -= vout.Value * common.Satoshi
					atomicalsID := atomicalsID(operation.RevealLocationTxID, int64(output_index))
					m.UTXOs[atomicalsID].AtomicalID = atomicalsID
					m.UTXOs[atomicalsID].DistributedFt = append(m.UTXOs[atomicalsID].DistributedFt, &UserDistributedInfo{
						UserPk:   tx.Vout[output_index].ScriptPubKey.Address,
						Name:     entity.Name,
						Nonce:    entity.Nonce,
						Time:     entity.Time,
						Bitworkc: entity.Bitworkc,
						Amount:   vout.Value * common.Satoshi,
					})
				}
			}
		}
	} else { // color_ft_atomicals_regular
		// a dmt in Vin has a total amount: entity.Amount,
		// color exactly the same amount of vout
		// burn the rest Dmt
		for _, vin := range tx.Vin {
			preAtomicalsID := atomicalsID(vin.Txid, int64(vin.Vout))
			if _, ok := m.UTXOs[preAtomicalsID]; !ok {
				continue
			}
			for index, entity := range m.UTXOs[preAtomicalsID].DistributedFt {
				remaining_value := entity.Amount
				m.UTXOs[preAtomicalsID].DistributedFt = append(m.UTXOs[preAtomicalsID].DistributedFt[:index], m.UTXOs[preAtomicalsID].DistributedFt[index+1:]...)
				for output_index, vout := range tx.Vout {
					if remaining_value < vout.Value*common.Satoshi { // if so, this Dmt will be burned
						break
					}
					remaining_value -= vout.Value * common.Satoshi
					atomicalsID := atomicalsID(operation.RevealLocationTxID, int64(output_index))
					m.UTXOs[atomicalsID].AtomicalID = atomicalsID
					m.UTXOs[atomicalsID].DistributedFt = append(m.UTXOs[atomicalsID].DistributedFt, &UserDistributedInfo{
						UserPk:   tx.Vout[output_index].ScriptPubKey.Address,
						Name:     entity.Name,
						Nonce:    entity.Nonce,
						Time:     entity.Time,
						Bitworkc: entity.Bitworkc,
						Amount:   vout.Value * common.Satoshi,
					})
				}
			}
		}
	}
	return nil
}
