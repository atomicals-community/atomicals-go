package atomicals

import (
	"strconv"

	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/atomicals-core/pkg/log"
	"github.com/btcsuite/btcd/btcjson"
)

// mintDistributedFt:operation dmt, Mint tokens of distributed mint type
func (m *Atomicals) mintDistributedFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	ticker := operation.Payload.Args.MintTicker
	ft, err := m.DistributedFtByName(ticker)
	if err != nil {
		log.Log.Panicf("DistributedFtByName err:%v", err)
	}
	if ft == nil {
		return errors.ErrNotDeployFt
	}
	ftEntity, err := m.DistributedFtByName(ticker)
	if err != nil {
		log.Log.Panicf("DistributedFtByName err:%v", err)
	}
	if ftEntity.CommitHeight <= operation.RevealLocationHeight-common.MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS {
		return errors.ErrInvalidCommitHeight
	}
	if operation.RevealLocationHeight < ftEntity.MintHeight {
		return errors.ErrInvalidMintHeight
	}
	if operation.CommitHeight < ftEntity.MintHeight {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsValidRevealLocationAndCommitVoutIndex() {
		return errors.ErrInvalidVinIndex
	}
	// if mint_amount == txout.value:
	if int64(vout[common.VOUT_EXPECT_OUTPUT_INDEX].Value*common.Satoshi) != ftEntity.MintAmount {
		return errors.ErrInvalidMintAmount
	}
	if ftEntity.MintMode == "perpetual" {
		if ftEntity.MaxMintsGlobal == ftEntity.MintedTimes {
			return nil
		}
		if ftEntity.Bci != "" {
			if operation.IsDftBitworkRolloverActivated() {
				success, _ := isTxidValidForPerpetualBitwork(operation.CommitTxID, ftEntity.Bv, ftEntity.MintedTimes, ftEntity.MaxMints, ftEntity.Bci, ftEntity.Bcs, true)
				if !success {
					return errors.ErrInvalidPerpetualBitwork
				}
			} else {
				success, _ := isTxidValidForPerpetualBitwork(operation.CommitTxID, ftEntity.Bv, ftEntity.MintedTimes, ftEntity.MaxMints, ftEntity.Bci, ftEntity.Bcs, false)
				if !success {
					return errors.ErrInvalidPerpetualBitwork
				}
			}
		}
		if ftEntity.Bri != "" {
			if operation.IsDftBitworkRolloverActivated() {
				success, _ := isTxidValidForPerpetualBitwork(operation.RevealLocationTxID, ftEntity.Bv, ftEntity.MintedTimes, ftEntity.MaxMints, ftEntity.Bri, ftEntity.Brs, true)
				if !success {
					return errors.ErrInvalidPerpetualBitwork
				}
			} else {
				success, _ := isTxidValidForPerpetualBitwork(operation.RevealLocationTxID, ftEntity.Bv, ftEntity.MintedTimes, ftEntity.MaxMints, ftEntity.Bri, ftEntity.Brs, false)
				if !success {
					return errors.ErrInvalidPerpetualBitwork
				}
			}
		}
	} else { //ftEntity.MintMode == "fixed"
		// ftEntity.MintBitworkc
		if ftEntity.MintedTimes > ftEntity.MaxMints {
			return errors.ErrInvalidMintedTimes
		} else if ftEntity.MintedTimes < ftEntity.MaxMints {
			_, _, err := witness.IsValidMintBitwork(operation.CommitTxID, operation.Payload.Args.MintBitworkc, operation.Payload.Args.MintBitworkr)
			if err != nil {
				return err
			}
		}
	}
	// if !operation.IsValidCommitVoutIndexForDmt() {
	// 	return errors.ErrInvalidVinIndex
	// }
	// if operation.RevealLocationHeight < common.ATOMICALS_ACTIVATION_HEIGHT_DMINT {
	// 	return errors.ErrInvalidRevealLocationHeight
	// }
	// if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
	// 	return errors.ErrInvalidCommitHeight
	// }
	// if !operation.IsWithinAcceptableBlocksForNameReveal() {
	// 	return errors.ErrInvalidCommitHeight
	// }
	// if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
	// 	return errors.ErrInvalidCommitHeight
	// }
	// if !operation.IsValidCommitVoutIndexForNameRevel() {
	// 	return errors.ErrInvalidCommitVoutIndex
	// }
	bitworkc, bitworkr, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
	locationID := operation.AtomicalsID
	entity := &db.UserFtInfo{
		MintTicker:  ticker,
		Nonce:       operation.Payload.Args.Nonce,
		Time:        operation.Payload.Args.Time,
		Bitworkc:    bitworkc,
		Bitworkr:    bitworkr,
		Amount:      int64(vout[common.VOUT_EXPECT_OUTPUT_INDEX].Value * common.Satoshi),
		AtomicalsID: locationID,
		LocationID:  locationID,
	}
	if err := m.InsertFtUTXO(entity); err != nil {
		log.Log.Panicf("InsertFtUTXO err:%v", err)
	}
	if err := m.UpdateDistributedFtAmount(ticker, ftEntity.MintedTimes+1); err != nil {
		log.Log.Panicf("UpdateDistributedFtAmount err:%v", err)
	}
	return nil
}

// is_txid_valid_for_perpetual_bitwork
func isTxidValidForPerpetualBitwork(txid string, bitwork_vec string, actual_mints, max_mints int64, mintBitworkrInc string, mintBitworkcStart int64, allow_higher bool) (bool, string) {
	starting_target := mintBitworkcStart
	target_increment, _ := strconv.Atoi(mintBitworkrInc) // never return err
	expected_minimum_bitwork := common.Calculate_expected_bitwork(bitwork_vec, actual_mints, max_mints, int64(target_increment), starting_target)
	if common.Is_mint_pow_valid(txid, expected_minimum_bitwork) {
		return true, expected_minimum_bitwork
	}
	if allow_higher {
		parts := common.ParseBitwork(expected_minimum_bitwork)
		if parts == nil {
			return false, ""
		}
		prefix := parts.Prefix
		next_full_bitwork_prefix := common.Get_next_bitwork_full_str(bitwork_vec, len(prefix))
		if common.Is_mint_pow_valid(txid, next_full_bitwork_prefix) {
			return true, next_full_bitwork_prefix
		}
	}
	return false, ""
}
