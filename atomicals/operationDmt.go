package atomicals

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

// mintFt: Mint tokens of distributed mint type (dft)
func (m *Atomicals) mintDistributedFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	atomicalsID := atomicalsID(operation.TxID, common.VOUT_EXPECT_OUTPUT_BYTES)
	ticker := operation.Payload.Args.MintTicker
	if _, ok := m.AtomicalsFtEntity[ticker]; !ok {
		return errors.ErrNotDeployFt
	}
	entity := &UserDmtEntity{
		Name:        ticker,
		Nonce:       operation.Payload.Args.Nonce,
		Time:        operation.Payload.Args.Time,
		Bitworkc:    operation.Payload.Args.Bitworkc,
		Amount:      vout[common.VOUT_EXPECT_OUTPUT_BYTES].Value * common.Satoshi,
		AtomicalsID: atomicalsID,
		Location:    atomicalsID,
	}
	ftEntity := m.AtomicalsFtEntity[ticker]
	if entity.Amount > ftEntity.MintAmount {
		return errors.ErrInvalidMintAmount
	}
	if operation.Height > ftEntity.MintHeight {
		return errors.ErrInvalidMintHeight
	}
	m.UTXOs[atomicalsID] = &AtomicalsUTXO{
		AtomicalID:    atomicalsID,
		DistributedFt: make([]*UserDmtEntity, 0),
	}
	m.UTXOs[atomicalsID].DistributedFt = append(m.UTXOs[atomicalsID].DistributedFt, entity)
	m.AtomicalsFtEntity[ticker].MintedAmount += entity.Amount
	return nil
}
