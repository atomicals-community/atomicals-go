package atomicals

import (
	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

// mintDirectFt: Mint fungible token with direct fixed supply
func (m *Atomicals) mintDirectFt(operation *witness.WitnessAtomicalsOperation, vout []btcjson.Vout, userPk string) (newGlobalDirectFt *postsql.GlobalDirectFt, err error) {
	if operation.RevealInputIndex != 0 {
		return nil, errors.ErrInvalidRevealInputIndex
	}
	if !operation.Payload.CheckRequest() {
		return nil, errors.ErrCheckRequest
	}
	if !utils.IsValidTicker(operation.Payload.Args.RequestTicker) {
		return nil, errors.ErrInvalidTicker
	}
	ft, err := m.DirectFtByName(operation.Payload.Args.RequestTicker)
	if err != nil {
		log.Log.Panicf("DistributedFtByName err:%v", err)
	}
	if ft != nil {
		return nil, errors.ErrTickerHasExist
	}
	if operation.Payload.IsImmutable() {
		return nil, errors.ErrCannotBeImmutable
	}
	if operation.Payload.Args.Bitworkc == "" {
		return nil, errors.ErrBitworkcNeeded
	}
	_, _, err = operation.IsValidBitwork()
	if err != nil {
		return nil, err
	}
	operation.CommitHeight, err = m.AtomicalsTxHeight(operation.CommitTxID)
	if err != nil {
		operation.CommitHeight, err = m.GetTxHeightByTxID(operation.CommitTxID)
		if err != nil {
			panic(err)
		}
	}
	if operation.CommitHeight < utils.ATOMICALS_ACTIVATION_HEIGHT {
		return nil, errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return nil, errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return nil, errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight > utils.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return nil, errors.ErrInvalidVinIndex
	}
	if operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return nil, errors.ErrInvalidVinIndex
	}

	amount := utils.MulSatoshi(vout[utils.VOUT_EXPECT_OUTPUT_INDEX].Value)
	newGlobalDirectFt = &postsql.GlobalDirectFt{
		UserPk:      userPk,
		AtomicalsID: operation.AtomicalsID,
		LocationID:  operation.LocationID,
		Type:        "FT",
		Subtype:     "direct",
		TickerName:  operation.Payload.Args.RequestTicker,
		// Meta:          operation.Payload.Meta,
		Bitworkc:  operation.Payload.Args.Bitworkc,
		Bitworkr:  operation.Payload.Args.Bitworkr,
		MaxSupply: amount,
	}
	return newGlobalDirectFt, nil
}
