package atomicals

import (
	"github.com/atomicals-go/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
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
	if !utils.IsValidTicker(operation.Payload.Args.RequestTicker) {
		return errors.ErrInvalidTicker
	}
	if operation.Payload.IsImmutable() {
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
		operation.CommitHeight, err = m.GetTxHeightByTxID(operation.CommitTxID)
		if err != nil {
			panic(err)
		}
	}
	if operation.CommitHeight < utils.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight > utils.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	if operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	entity := &postsql.GlobalDirectFt{
		UserPk:      userPk,
		AtomicalsID: operation.AtomicalsID,
		LocationID:  operation.LocationID,
		Type:        "FT",
		Subtype:     "direct",
		TickerName:  operation.Payload.Args.RequestTicker,
		// Meta:          operation.Payload.Meta,
		Bitworkc:  operation.Payload.Args.Bitworkc,
		Bitworkr:  operation.Payload.Args.Bitworkr,
		MaxSupply: int64(vout[utils.VOUT_EXPECT_OUTPUT_INDEX].Value * utils.Satoshi),
	}
	if err := m.InsertDirectFtUTXO(entity); err != nil {
		log.Log.Panicf("InsertDirectFtUTXO err:%v", err)
	}
	return nil
}
