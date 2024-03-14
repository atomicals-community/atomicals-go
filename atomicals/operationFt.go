package atomicals

import (
	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/atomicals-core/pkg/log"
	"github.com/btcsuite/btcd/btcjson"
)

// mintDirectFt: Mint fungible token with direct fixed supply
func (m *Atomicals) mintDirectFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	if operation.CommitHeight >= operation.RevealLocationHeight-common.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS {
		return errors.ErrInvalidCommitHeight
	}
	if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if operation.CommitVoutIndex != common.VOUT_EXPECT_OUTPUT_INDEX {
		return errors.ErrInvalidVinIndex
	}
	if operation.IsImmutable() {
		return errors.ErrCannotBeImmutable
	}
	bitworkc, bitworkr, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
	locationID := operation.AtomicalsID
	atomicalsFtInfo := &db.UserFtInfo{
		UserPk:        userPk,
		AtomicalsID:   locationID,
		LocationID:    locationID,
		Type:          "FT",
		Subtype:       "direct",
		RequestTicker: operation.Payload.Args.RequestTicker,
		Meta:          operation.Payload.Meta,
		Bitworkc:      bitworkc,
		Bitworkr:      bitworkr,
		MaxSupply:     int64(vout[common.VOUT_EXPECT_OUTPUT_INDEX].Value * common.Satoshi),
	}
	if !common.IsValidTicker(atomicalsFtInfo.RequestTicker) {
		return errors.ErrInvalidTicker
	}
	if err := m.InsertFtUTXO(atomicalsFtInfo); err != nil {
		log.Log.Panicf("InsertFtUTXO err:%v", err)
	}
	return nil
}
