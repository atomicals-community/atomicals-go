package atomicals

import (
	"strconv"

	"github.com/atomicals-core/atomicals/DB/postsql"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/atomicals-core/pkg/log"
)

// deployDistributedFt: operation dft
func (m *Atomicals) deployDistributedFt(operation *witness.WitnessAtomicalsOperation, userPk string) error {
	if operation.RevealInputIndex != 0 {
		return errors.ErrInvalidRevealInputIndex
	}
	if !operation.Payload.CheckRequest() {
		return errors.ErrCheckRequest
	}
	if !common.IsValidTicker(operation.Payload.Args.RequestTicker) {
		return errors.ErrInvalidTicker
	}
	ft, err := m.DistributedFtByName(operation.Payload.Args.RequestTicker)
	if err != nil {
		log.Log.Panicf("DistributedFtByName err:%v", err)
	}
	if ft != nil {
		return errors.ErrTickerHasExist
	}
	if operation.Payload.Args.MintHeight < common.DFT_MINT_HEIGHT_MIN || common.DFT_MINT_HEIGHT_MAX < operation.Payload.Args.MintHeight {
		return errors.ErrInvalidMintHeight
	}
	if operation.Payload.Args.MintAmount < common.DFT_MINT_AMOUNT_MIN || common.DFT_MINT_AMOUNT_MAX < operation.Payload.Args.MintAmount {
		return errors.ErrInvalidMintHeight
	}
	if operation.Payload.Args.MaxMints < common.DFT_MINT_MAX_MIN_COUNT {
		return errors.ErrInvalidMaxMints
	}
	if operation.RevealLocationHeight < common.ATOMICALS_ACTIVATION_HEIGHT_DENSITY {
		if operation.Payload.Args.MaxMints > common.DFT_MINT_MAX_MAX_COUNT_LEGACY {
			return errors.ErrInvalidMaxMints
		}
	} else {
		if operation.Payload.Args.MaxMints > common.DFT_MINT_MAX_MAX_COUNT_DENSITY {
			return errors.ErrInvalidMaxMints
		}
	}
	mintBitworkc, _, err := witness.IsValidMintBitwork(operation.CommitTxID, operation.Payload.Args.MintBitworkc, operation.Payload.Args.MintBitworkr)
	if err != nil {
		return err
	}
	if mintBitworkc != nil && len(mintBitworkc.Prefix) < 4 {
		return errors.ErrInvalidBitworkcPrefix
	}
	if operation.IsImmutable() {
		return errors.ErrCannotBeImmutable
	}
	if operation.Payload.Args.Md != "" && operation.Payload.Args.Md != "0" && operation.Payload.Args.Md != "1" {
		return errors.ErrInvalidDftMd
	}
	if operation.Payload.Args.Bitworkc == "" {
		return errors.ErrBitworkcNeeded
	}
	_, _, err = operation.IsValidBitwork()
	if err != nil {
		return err
	}
	operation.CommitHeight, err = m.GetCommitHeight(operation.CommitTxID)
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
	atomicalsID := operation.AtomicalsID
	entity := &postsql.GlobalDistributedFt{
		AtomicalsID:  atomicalsID,
		TickerName:   operation.Payload.Args.RequestTicker,
		Type:         "FT",
		Subtype:      "decentralized",
		MintAmount:   operation.Payload.Args.MintAmount,
		MaxMints:     operation.Payload.Args.MaxMints,
		MintHeight:   operation.Payload.Args.MintHeight,
		MintBitworkc: operation.Payload.Args.MintBitworkc,
		MintBitworkr: operation.Payload.Args.MintBitworkr,
		Bitworkc:     operation.Payload.Args.Bitworkc,
		Bitworkr:     operation.Payload.Args.Bitworkr,
		// Meta:         operation.Payload.Meta,
		MintedTimes:  0,
		Md:           operation.Payload.Args.Md,
		Bv:           operation.Payload.Args.Bv,
		Bci:          operation.Payload.Args.Bci,
		Bri:          operation.Payload.Args.Bri,
		Bcs:          operation.Payload.Args.Bcs,
		Brs:          operation.Payload.Args.Brs,
		Maxg:         operation.Payload.Args.Maxg,
		CommitHeight: operation.CommitHeight,
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
		if entity.MaxMintsGlobal != 0 {
			entity.MaxSupply = entity.MintAmount * entity.MaxMintsGlobal
		} else {
			entity.MaxSupply = -1
		}
	} else {
		entity.MintMode = "fixed"
		entity.MaxSupply = entity.MintAmount * entity.MaxMints
	}
	if err := m.InsertDistributedFt(entity); err != nil {
		log.Log.Panicf("InsertDistributedFt err:%v", err)
	}
	return nil
}
