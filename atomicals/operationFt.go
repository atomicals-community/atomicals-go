package atomicals

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

// mintDirectFt: Mint fungible token with direct fixed supply
func (m *Atomicals) mintDirectFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	atomicalsID := atomicalsID(operation.TxID, common.VOUT_EXPECT_OUTPUT_BYTES)
	atomicalsFtInfo := &UserFtEntity{
		UserPk:      userPk,
		CommitTxID:  vin.Txid,
		CommitIndex: int64(vin.Vout),
		// CommitHeight:    ,
		CurrentTxID:   operation.TxID,
		CurrentHeight: operation.Height,
		Type:          "FT",
		Subtype:       "direct",
		RequestTicker: operation.Payload.Args.RequestTicker,
		Meta:          operation.Payload.Meta,
		Bitworkc:      operation.Payload.Args.Bitworkc,
		MaxSupply:     int64(vout[common.VOUT_EXPECT_OUTPUT_BYTES].Value * common.Satoshi),
	}
	if atomicalsFtInfo.CommitHeight >= atomicalsFtInfo.CurrentHeight-common.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS {
		return errors.ErrInvalidFtCommitHeight
	}
	if atomicalsFtInfo.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidFtCurrentHeight
	}
	if atomicalsFtInfo.CurrentHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidFtCurrentHeight
	}
	if !common.IsValidTicker(atomicalsFtInfo.RequestTicker) {
		return errors.ErrInvalidTicker
	}
	// # Check if there was requested proof of work, && if there was then only allow the mint to happen if it was successfully executed the proof of work
	bitwork := common.ParseBitwork(atomicalsFtInfo.Bitworkc)
	if bitwork != nil {
		if !common.IsProofOfWorkPrefixMatch(atomicalsFtInfo.CommitTxID, bitwork.Prefix, bitwork.Ext) {
			return errors.ErrInvalidBitWork
		}
	}
	m.UTXOs[atomicalsID] = &AtomicalsUTXO{
		AtomicalID: atomicalsID,
		DirectFt:   make([]*UserFtEntity, 0),
	}
	m.UTXOs[atomicalsID].DirectFt = append(m.UTXOs[atomicalsID].DirectFt, atomicalsFtInfo)
	return nil
}
