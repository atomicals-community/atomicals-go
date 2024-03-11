package atomicals

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

// mintFt: Mint tokens of distributed mint type (dft)
func (m *Atomicals) mintDistributedFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	if operation.RevealLocationHeight < common.ATOMICALS_ACTIVATION_HEIGHT_DMINT {
		return errors.ErrInvalidRevealLocationHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsValidCommitVoutIndexForNameRevel() {
		return errors.ErrInvalidCommitVoutIndex
	}
	bitworkc, bitworkr, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
	atomicalsID := atomicalsID(operation.RevealLocationTxID, operation.RevealLocationVoutIndex)
	ticker := operation.Payload.Args.MintTicker
	if !m.DistributedFtHasExist(ticker) {
		return errors.ErrNotDeployFt
	}
	entity := &UserDistributedInfo{
		Name:        ticker,
		Nonce:       operation.Payload.Args.Nonce,
		Time:        operation.Payload.Args.Time,
		Bitworkc:    bitworkc,
		Bitworkr:    bitworkr,
		Amount:      vout[common.VOUT_EXPECT_OUTPUT_INDEX].Value * common.Satoshi,
		AtomicalsID: atomicalsID,
		Location:    atomicalsID,
	}
	ftEntity := m.GlobalDistributedFtMap[ticker]
	if entity.Amount > ftEntity.MintAmount {
		return errors.ErrInvalidMintAmount
	}
	if ftEntity.MaxMints < ftEntity.MintedAmount+entity.Amount {
		return errors.ErrInvalidMintHeight
	}
	if operation.RevealLocationHeight > ftEntity.MintHeight {
		return errors.ErrInvalidMintHeight
	}
	m.ensureUTXONotNil(atomicalsID)
	m.UTXOs[atomicalsID].DistributedFt = append(m.UTXOs[atomicalsID].DistributedFt, entity)
	m.GlobalDistributedFtMap[ticker].MintedAmount += entity.Amount
	return nil
}
