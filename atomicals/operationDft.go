package atomicals

import (
	"strconv"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) deployDistributedFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	if operation.CommitHeight >= operation.RevealLocationHeight-common.MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS {
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
	mintBitworkc, mintBitworkr, err := witness.IsValidMintBitwork(operation.CommitTxID, operation.Payload.Args.MintBitworkc, operation.Payload.Args.MintBitworkr)
	if err != nil {
		return err
	}
	atomicalsID := atomicalsID(operation.RevealLocationTxID, operation.RevealLocationVoutIndex)
	entity := &DistributedFtInfo{
		AtomicalsID:  atomicalsID,
		Ticker:       operation.Payload.Args.RequestTicker,
		Type:         "FT",
		Subtype:      "decentralized",
		MintAmount:   operation.Payload.Args.MintAmount,
		MaxMints:     operation.Payload.Args.MaxMints,
		MintHeight:   operation.Payload.Args.MintHeight,
		MintBitworkc: mintBitworkc,
		MintBitworkr: mintBitworkr,
		Bitworkc:     bitworkc,
		Bitworkr:     bitworkr,
		Meta:         operation.Payload.Meta,
		MintedTimes:  0,
		Md:           operation.Payload.Args.Md,
		Bv:           operation.Payload.Args.Bv,
		Bci:          operation.Payload.Args.Bci,
		Bri:          operation.Payload.Args.Bri,
		Bcs:          operation.Payload.Args.Bcs,
		Brs:          operation.Payload.Args.Brs,
		Maxg:         operation.Payload.Args.Maxg,
	}
	if entity.Bitworkc != nil && len(entity.Bitworkc.Prefix) < 4 {
		return errors.ErrInvalidBitworkcPrefix
	}
	if entity.Bitworkc == nil {
		return errors.ErrNameTypeMintMastHaveBitworkc
	}
	if !common.IsValidTicker(entity.Ticker) {
		return errors.ErrInvalidTicker
	}
	if m.DistributedFtHasExist(entity.Ticker) {
		return errors.ErrTickerHasExist
	}
	if common.DFT_MINT_HEIGHT_MAX < entity.MintHeight {
		return errors.ErrInvalidMintHeight
	}
	if entity.MintAmount < common.DFT_MINT_AMOUNT_MIN || common.DFT_MINT_AMOUNT_MAX < entity.MintAmount {
		return errors.ErrInvalidMintHeight
	}
	if entity.MaxMints < common.DFT_MINT_MAX_MIN_COUNT {
		return errors.ErrInvalidMaxMints
	}
	if operation.RevealLocationHeight < common.ATOMICALS_ACTIVATION_HEIGHT_DENSITY {
		if entity.MaxMints > common.DFT_MINT_MAX_MAX_COUNT_LEGACY {
			return errors.ErrInvalidMaxMints
		}
	} else {
		if entity.MaxMints > common.DFT_MINT_MAX_MAX_COUNT_DENSITY {
			return errors.ErrInvalidMaxMints
		}
	}
	if entity.Md != "" && entity.Md != "0" && entity.Md != "1" {
		return errors.ErrInvalidDftMd
	}
	if common.ATOMICALS_ACTIVATION_HEIGHT_DENSITY <= operation.RevealLocationHeight && entity.Md == "1" {
		if !common.IsHexStringRegex(operation.Payload.Args.Bv) || len(operation.Payload.Args.Bv) < 4 {
			return errors.ErrInvalidDftBv
		}
		if operation.Payload.Args.MintBitworkc != "" || operation.Payload.Args.MintBitworkr != "" {
			return errors.ErrInvalidDftMintBitwork
		}
		if operation.Payload.Args.Bci != "" {
			bci, err := strconv.Atoi(operation.Payload.Args.Bci)
			if err == nil {
				if 64 < bci {
					return errors.ErrInvalidDftBci
				}
				if operation.Payload.Args.Bcs < 64 || 256 < operation.Payload.Args.Bcs {
					return errors.ErrInvalidDftBsc
				}
			}
		}
		if operation.Payload.Args.Bri != "" {
			bri, err := strconv.Atoi(operation.Payload.Args.Bri)
			if err == nil {
				if 64 < bri {
					return errors.ErrInvalidDftBri
				}
				if operation.Payload.Args.Brs < 64 || 256 < operation.Payload.Args.Brs {
					return errors.ErrInvalidDftBrs
				}
			}
		}
		if 100000 < operation.Payload.Args.MaxMints {
			return errors.ErrInvalidMaxMints
		}
		if operation.Payload.Args.Maxg < common.DFT_MINT_MAX_MIN_COUNT || common.DFT_MINT_MAX_MAX_COUNT_DENSITY < operation.Payload.Args.Maxg {
			return errors.ErrInvalidDftMaxg
		}
		entity.MaxMintsGlobal = operation.Payload.Args.Maxg
		entity.MintMode = "perpetual"
		entity.MaxSupply = entity.MintAmount * entity.MaxMintsGlobal
	} else {
		entity.MintMode = "fixed"
		entity.MaxSupply = -1
	}
	m.GlobalDistributedFtMap[entity.Ticker] = entity
	return nil

}
