package atomicals

import (
	"strconv"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
)

// deployDistributedFt: operation dft
func (m *Atomicals) deployDistributedFt(operation *witness.WitnessAtomicalsOperation, userPk string) (newGlobalDistributedFt *postsql.GlobalDistributedFt, err error) {
	if operation.RevealInputIndex != 0 {
		return nil, errors.ErrInvalidRevealInputIndex
	}
	if !operation.Payload.CheckRequest() {
		return nil, errors.ErrCheckRequest
	}
	if !utils.IsValidTicker(operation.Payload.Args.RequestTicker) {
		return nil, errors.ErrInvalidTicker
	}
	ft, err := m.DistributedFtByName(operation.Payload.Args.RequestTicker)
	if err != nil {
		log.Log.Panicf("DistributedFtByName err:%v", err)
	}
	if ft != nil {
		return nil, errors.ErrTickerHasExist
	}
	if operation.Payload.Args.Bitworkc == "" {
		return nil, errors.ErrBitworkcNeeded
	}
	_, _, err = operation.IsValidBitwork()
	if err != nil {
		return nil, err
	}
	if operation.Payload.Args.MintHeight < utils.DFT_MINT_HEIGHT_MIN || utils.DFT_MINT_HEIGHT_MAX < operation.Payload.Args.MintHeight {
		return nil, errors.ErrInvalidMintHeight
	}
	if operation.Payload.Args.MintAmount < utils.DFT_MINT_AMOUNT_MIN || utils.DFT_MINT_AMOUNT_MAX < operation.Payload.Args.MintAmount {
		return nil, errors.ErrInvalidMintHeight
	}
	if operation.Payload.Args.MaxMints < utils.DFT_MINT_MAX_MIN_COUNT {
		return nil, errors.ErrInvalidMaxMints
	}
	if operation.RevealLocationHeight < utils.ATOMICALS_ACTIVATION_HEIGHT_DENSITY {
		if operation.Payload.Args.MaxMints > utils.DFT_MINT_MAX_MAX_COUNT_LEGACY {
			return nil, errors.ErrInvalidMaxMints
		}
	} else {
		if operation.Payload.Args.MaxMints > utils.DFT_MINT_MAX_MAX_COUNT_DENSITY {
			return nil, errors.ErrInvalidMaxMints
		}
	}
	mintBitworkc, _, err := utils.ParseMintBitwork(operation.CommitTxID, operation.Payload.Args.MintBitworkc, operation.Payload.Args.MintBitworkr)
	if err != nil {
		return nil, err
	}
	if mintBitworkc != nil && len(mintBitworkc.Prefix) < 4 {
		return nil, errors.ErrInvalidBitworkcPrefix
	}
	if operation.Payload.IsImmutable() {
		return nil, errors.ErrCannotBeImmutable
	}
	if operation.Payload.Args.Md != "" && operation.Payload.Args.Md != "0" && operation.Payload.Args.Md != "1" {
		return nil, errors.ErrInvalidDftMd
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
	newGlobalDistributedFt = &postsql.GlobalDistributedFt{
		AtomicalsID:  operation.AtomicalsID,
		LocationID:   operation.LocationID,
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

	if utils.ATOMICALS_ACTIVATION_HEIGHT_DENSITY <= operation.RevealLocationHeight && newGlobalDistributedFt.Md == "1" {
		if !utils.IsHexStringRegex(operation.Payload.Args.Bv) || len(operation.Payload.Args.Bv) < 4 {
			return nil, errors.ErrInvalidDftBv
		}
		if operation.Payload.Args.MintBitworkc != "" || operation.Payload.Args.MintBitworkr != "" {
			return nil, errors.ErrInvalidDftMintBitwork
		}
		if operation.Payload.Args.Bci != "" {
			bci, err := strconv.Atoi(operation.Payload.Args.Bci)
			if err == nil {
				if 64 < bci {
					return nil, errors.ErrInvalidDftBci
				}
				if operation.Payload.Args.Bcs < 64 || 256 < operation.Payload.Args.Bcs {
					return nil, errors.ErrInvalidDftBsc
				}
			}
		}
		if operation.Payload.Args.Bri != "" {
			bri, err := strconv.Atoi(operation.Payload.Args.Bri)
			if err == nil {
				if 64 < bri {
					return nil, errors.ErrInvalidDftBri
				}
				if operation.Payload.Args.Brs < 64 || 256 < operation.Payload.Args.Brs {
					return nil, errors.ErrInvalidDftBrs
				}
			}
		}
		if 100000 < operation.Payload.Args.MaxMints {
			return nil, errors.ErrInvalidMaxMints
		}
		if operation.Payload.Args.Maxg < utils.DFT_MINT_MAX_MIN_COUNT || utils.DFT_MINT_MAX_MAX_COUNT_DENSITY < operation.Payload.Args.Maxg {
			return nil, errors.ErrInvalidDftMaxg
		}
		newGlobalDistributedFt.MaxMintsGlobal = operation.Payload.Args.Maxg
		newGlobalDistributedFt.MintMode = "perpetual"
		if newGlobalDistributedFt.MaxMintsGlobal != 0 {
			newGlobalDistributedFt.MaxSupply = newGlobalDistributedFt.MintAmount * newGlobalDistributedFt.MaxMintsGlobal
		} else {
			newGlobalDistributedFt.MaxSupply = -1
		}
	} else {
		newGlobalDistributedFt.MintMode = "fixed"
		newGlobalDistributedFt.MaxSupply = newGlobalDistributedFt.MintAmount * newGlobalDistributedFt.MaxMints
	}
	return newGlobalDistributedFt, nil
}
