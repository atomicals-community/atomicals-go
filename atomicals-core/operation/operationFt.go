package atomicals

import (
	"github.com/atomicals-go/atomicals-core/common"
	"github.com/atomicals-go/atomicals-core/repo/postsql"
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/btcsuite/btcd/btcjson"
)

// mintDirectFt: Mint fungible token with direct fixed supply
func (m *Atomicals) mintDirectFt(operation *witness.WitnessAtomicalsOperation, vout []btcjson.Vout, userPk string) (err error) {
	if operation.RevealInputIndex != 0 {
		return errors.ErrInvalidRevealInputIndex
	}
	if !operation.Payload.CheckRequest() {
		return errors.ErrCheckRequest
	}
	if !common.IsValidTicker(operation.Payload.Args.RequestTicker) {
		return errors.ErrInvalidTicker
	}
	if operation.IsImmutable() {
		return errors.ErrCannotBeImmutable
	}
	if operation.Payload.Args.Bitworkc == "" {
		return errors.ErrBitworkcNeeded
	}
	_, _, err = operation.IsValidBitwork()
	if err != nil {
		return err
	}
	operation.CommitHeight, err = m.BtcTxHeight(operation.CommitTxID)
	if err != nil {
		panic(err)
	}
	if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight > common.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != common.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	if operation.CommitVoutIndex != common.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}

	locationID := operation.AtomicalsID
	atomicalsFtInfo := &postsql.UTXOFtInfo{
		UserPk:        userPk,
		AtomicalsID:   locationID,
		LocationID:    locationID,
		Type:          "FT",
		Subtype:       "direct",
		RequestTicker: operation.Payload.Args.RequestTicker,
		// Meta:          operation.Payload.Meta,
		Bitworkc:  operation.Payload.Args.Bitworkc,
		Bitworkr:  operation.Payload.Args.Bitworkr,
		MaxSupply: int64(vout[common.VOUT_EXPECT_OUTPUT_INDEX].Value * common.Satoshi),
	}
	if err := m.InsertFtUTXO(atomicalsFtInfo); err != nil {
		log.Log.Panicf("InsertFtUTXO err:%v", err)
	}
	return nil
}
