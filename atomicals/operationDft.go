package atomicals

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) deployFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	atomicalsID := atomicalsID(operation.TxID, common.VOUT_EXPECT_OUTPUT_BYTES)
	atomicalsDftInfo := &AtomicalsFtInfo{
		AtomicalsID:  atomicalsID,
		Ticker:       operation.Payload.Args.RequestTicker,
		MintAmount:   operation.Payload.Args.MintAmount,
		MaxMints:     operation.Payload.Args.MaxMints,
		MintHeight:   operation.Payload.Args.MintHeight,
		Bitworkc:     operation.Payload.Args.Bitworkc,
		Meta:         operation.Payload.Meta,
		MintedAmount: 0,
	}
	if !common.IsValidTicker(atomicalsDftInfo.Ticker) {
		return errors.ErrInvalidTicker
	}
	if _, ok := m.AtomicalsFtEntity[atomicalsDftInfo.Ticker]; ok {
		return errors.ErrTickerHasExist
	}
	if operation.Height < common.ATOMICALS_ACTIVATION_HEIGHT_DENSITY {
		if atomicalsDftInfo.MaxMints > common.DFT_MINT_MAX_MAX_COUNT_LEGACY {
			return errors.ErrInvalidMaxMints
		}
	} else {
		if atomicalsDftInfo.MaxMints > common.DFT_MINT_MAX_MAX_COUNT_DENSITY {
			return errors.ErrInvalidMaxMints
		}
	}
	bitwork := common.ParseBitwork(atomicalsDftInfo.Bitworkc)
	if bitwork != nil {
		if !common.IsProofOfWorkPrefixMatch(operation.TxID, bitwork.Prefix, bitwork.Ext) {
			return errors.ErrInvalidBitWork
		}
	}
	m.AtomicalsFtEntity[atomicalsDftInfo.Ticker] = atomicalsDftInfo
	return nil
}
