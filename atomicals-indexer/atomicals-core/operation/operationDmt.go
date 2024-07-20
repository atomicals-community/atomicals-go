package atomicals

import (
	"strconv"

	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

// mintDistributedFt:operation dmt, Mint tokens of distributed mint type
func (m *Atomicals) mintDistributedFt(operation *witness.WitnessAtomicalsOperation, vout []btcjson.Vout, userPk string) (newUTXOFtInfo *postsql.UTXOFtInfo, updateDistributedFt *postsql.GlobalDistributedFt, err error) {
	if operation.RevealInputIndex != 0 {
		return nil, nil, errors.ErrInvalidRevealInputIndex
	}
	ticker := operation.Payload.Args.MintTicker
	updateDistributedFt, err = m.DistributedFtByName(ticker)
	if err != nil {
		log.Log.Panicf("DistributedFtByName err:%v", err)
	}
	if updateDistributedFt == nil {
		return nil, nil, errors.ErrNotDeployFt
	}
	if operation.RevealLocationHeight < updateDistributedFt.CommitHeight+utils.MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS {
		return nil, nil, errors.ErrInvalidCommitHeight
	}
	operation.CommitHeight, err = m.AtomicalsTxHeight(operation.CommitTxID)
	if err != nil {
		operation.CommitHeight, err = m.GetTxHeightByTxID(operation.CommitTxID)
		if err != nil {
			panic(err)
		}
	}
	if operation.CommitHeight < updateDistributedFt.MintHeight {
		return nil, nil, errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight >= utils.ATOMICALS_ACTIVATION_HEIGHT_COMMITZ && operation.CommitVoutIndex != utils.VOUT_EXPECT_OUTPUT_INDEX {
		return nil, nil, errors.ErrInvalidVinIndex
	}
	// if mint_amount == txout.value:
	amount := utils.MulSatoshi(vout[utils.VOUT_EXPECT_OUTPUT_INDEX].Value)
	if amount != updateDistributedFt.MintAmount {
		return nil, nil, errors.ErrInvalidMintAmount
	}

	if updateDistributedFt.MintMode == "perpetual" {
		if updateDistributedFt.MaxMintsGlobal == updateDistributedFt.MintedTimes {
			return nil, nil, nil
		}
		if updateDistributedFt.Bci != "" {
			if operation.IsDftBitworkRolloverActivated() {
				success, _ := isTxidValidForPerpetualBitwork(operation.CommitTxID, updateDistributedFt.Bv, updateDistributedFt.MintedTimes, updateDistributedFt.MaxMints, updateDistributedFt.Bci, updateDistributedFt.Bcs, true)
				if !success {
					return nil, nil, errors.ErrInvalidPerpetualBitwork
				}
			} else {
				success, _ := isTxidValidForPerpetualBitwork(operation.CommitTxID, updateDistributedFt.Bv, updateDistributedFt.MintedTimes, updateDistributedFt.MaxMints, updateDistributedFt.Bci, updateDistributedFt.Bcs, false)
				if !success {
					return nil, nil, errors.ErrInvalidPerpetualBitwork
				}
			}
		}
		if updateDistributedFt.Bri != "" {
			if operation.IsDftBitworkRolloverActivated() {
				success, _ := isTxidValidForPerpetualBitwork(operation.RevealLocationTxID, updateDistributedFt.Bv, updateDistributedFt.MintedTimes, updateDistributedFt.MaxMints, updateDistributedFt.Bri, updateDistributedFt.Brs, true)
				if !success {
					return nil, nil, errors.ErrInvalidPerpetualBitwork
				}
			} else {
				success, _ := isTxidValidForPerpetualBitwork(operation.RevealLocationTxID, updateDistributedFt.Bv, updateDistributedFt.MintedTimes, updateDistributedFt.MaxMints, updateDistributedFt.Bri, updateDistributedFt.Brs, false)
				if !success {
					return nil, nil, errors.ErrInvalidPerpetualBitwork
				}
			}
		}
	} else { //updateDistributedFt.MintMode == "fixed"
		if updateDistributedFt.MintedTimes > updateDistributedFt.MaxMints {
			return nil, nil, errors.ErrInvalidMintedTimes
		} else if updateDistributedFt.MintedTimes < updateDistributedFt.MaxMints {
			bitworkc, bitworkr, err := utils.ParseMintBitwork(operation.CommitTxID, operation.Payload.Args.MintBitworkc, operation.Payload.Args.MintBitworkr)
			if err != nil {
				return nil, nil, err
			}
			if bitworkc != nil {
				if !utils.IsProofOfWorkPrefixMatch(operation.CommitTxID, bitworkc.Prefix, bitworkc.Ext) {
					return nil, nil, errors.ErrInvalidBitWork
				}
			}
			if bitworkr != nil {
				if !utils.IsProofOfWorkPrefixMatch(operation.CommitTxID, bitworkr.Prefix, bitworkr.Ext) {
					return nil, nil, errors.ErrInvalidBitWork
				}
			}
		}
	}
	_, _, err = operation.IsValidBitwork()
	if err != nil {
		return nil, nil, err
	}
	newUTXOFtInfo = &postsql.UTXOFtInfo{
		UserPk:      userPk,
		MintTicker:  ticker,
		Nonce:       operation.Payload.Args.Nonce,
		Time:        operation.Payload.Args.Time,
		Bitworkc:    operation.Payload.Args.Bitworkc,
		Bitworkr:    operation.Payload.Args.Bitworkr,
		Amount:      amount,
		AtomicalsID: updateDistributedFt.AtomicalsID,
		LocationID:  operation.LocationID,
	}
	updateDistributedFt.MintedTimes = updateDistributedFt.MintedTimes + 1
	return newUTXOFtInfo, updateDistributedFt, nil
}

// is_txid_valid_for_perpetual_bitwork
func isTxidValidForPerpetualBitwork(txid string, bitworkVec string, actualMints, maxMints int64, mintBitworkrInc string, mintBitworkcStart int64, allowHigher bool) (bool, string) {
	startingTarget := mintBitworkcStart
	targetIncrement, _ := strconv.Atoi(mintBitworkrInc) // never return err
	expectedMinimumBitwork := utils.Calculate_expected_bitwork(bitworkVec, actualMints, maxMints, int64(targetIncrement), startingTarget)
	if utils.IsMintPowValid(txid, expectedMinimumBitwork) {
		return true, expectedMinimumBitwork
	}
	if allowHigher {
		parts := utils.ParseBitwork(expectedMinimumBitwork)
		if parts == nil {
			return false, ""
		}
		prefix := parts.Prefix
		nextFullBitworkPrefix := utils.GetNextBitworkFullStr(bitworkVec, len(prefix))
		if utils.IsMintPowValid(txid, nextFullBitworkPrefix) {
			return true, nextFullBitworkPrefix
		}
	}
	return false, ""
}
