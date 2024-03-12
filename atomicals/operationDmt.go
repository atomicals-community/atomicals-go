package atomicals

import (
	"strconv"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/atomicals/witness"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
)

// mintFt: Mint tokens of distributed mint type (dft)
func (m *Atomicals) mintDistributedFt(operation *witness.WitnessAtomicalsOperation, vin btcjson.Vin, vout []btcjson.Vout, userPk string) error {
	if !operation.IsValidCommitVoutIndexForDmt() {
		return errors.ErrInvalidVinIndex
	}
	if operation.RevealLocationHeight < common.ATOMICALS_ACTIVATION_HEIGHT_DMINT {
		return errors.ErrInvalidRevealLocationHeight
	}
	if !operation.IsWithinAcceptableBlocksForGeneralReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsWithinAcceptableBlocksForNameReveal() {
		return errors.ErrInvalidCommitHeight
	}
	if operation.CommitHeight < common.ATOMICALS_ACTIVATION_HEIGHT {
		return errors.ErrInvalidCommitHeight
	}
	if !operation.IsValidCommitVoutIndexForNameRevel() {
		return errors.ErrInvalidCommitVoutIndex
	}
	bitworkc, bitworkr, err := operation.IsValidBitwork()
	if err != nil {
		return err
	}
	atomicalsID := atomicalsID(operation.RevealLocationTxID, operation.RevealLocationVoutIndex)
	ticker := operation.Payload.Args.MintTicker
	if !m.DistributedFtHasExist(ticker) {
		return errors.ErrNotDeployFt
	}
	entity := &UserDistributedInfo{
		Name:        ticker,
		Nonce:       operation.Payload.Args.Nonce,
		Time:        operation.Payload.Args.Time,
		Bitworkc:    bitworkc,
		Bitworkr:    bitworkr,
		Amount:      int64(vout[common.VOUT_EXPECT_OUTPUT_INDEX].Value * common.Satoshi),
		AtomicalsID: atomicalsID,
		Location:    atomicalsID,
	}
	ftEntity := m.GlobalDistributedFtMap[ticker]
	if operation.RevealLocationHeight < ftEntity.MintHeight {
		return errors.ErrInvalidMintHeight
	}
	if operation.CommitHeight < ftEntity.MintHeight {
		return errors.ErrInvalidCommitHeight
	}
	if entity.Amount != ftEntity.MintAmount {
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
	m.ensureUTXONotNil(atomicalsID)
	m.UTXOs[atomicalsID].DistributedFt = append(m.UTXOs[atomicalsID].DistributedFt, entity)
	m.GlobalDistributedFtMap[ticker].MintedTimes += 1
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
