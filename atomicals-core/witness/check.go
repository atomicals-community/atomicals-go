package witness

import (
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/utils"
)

// is_dft_bitwork_rollover_activated
func (m *WitnessAtomicalsOperation) IsDftBitworkRolloverActivated() bool {
	return m.RevealLocationHeight >= utils.ATOMICALS_ACTIVATION_HEIGHT_DFT_BITWORK_ROLLOVER
}

// is_within_acceptable_blocks_for_name_reveal
func (m *WitnessAtomicalsOperation) IsWithinAcceptableBlocksForNameReveal() bool {
	return m.CommitHeight >= m.RevealLocationHeight-utils.MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS
}

// is_within_acceptable_blocks_for_general_reveal
func (m *WitnessAtomicalsOperation) IsWithinAcceptableBlocksForGeneralReveal() bool {
	return m.CommitHeight >= m.RevealLocationHeight-utils.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS
}

func (m *WitnessAtomicalsOperation) IsValidBitwork() (*utils.Bitwork, *utils.Bitwork, error) {
	if m.Payload == nil {
		return nil, nil, nil
	}
	if m.Payload.Args == nil {
		return nil, nil, nil
	}
	bitworkc := utils.ParseBitwork(m.Payload.Args.Bitworkc)
	if bitworkc != nil {
		if !utils.IsProofOfWorkPrefixMatch(m.CommitTxID, bitworkc.Prefix, bitworkc.Ext) {
			return nil, nil, errors.ErrInvalidBitWork
		}
	}
	bitworkr := utils.ParseBitwork(m.Payload.Args.Bitworkr)
	if bitworkr != nil {
		if !utils.IsProofOfWorkPrefixMatch(m.CommitTxID, bitworkr.Prefix, bitworkr.Ext) {
			return nil, nil, errors.ErrInvalidBitWork
		}
	}
	return bitworkc, bitworkr, nil
}

// is_splat_operation
func (m *WitnessAtomicalsOperation) IsSplatOperation() bool {
	return m != nil && m.Op == "x" && m.RevealInputIndex == 0
}

// is_split_operation
func (m *WitnessAtomicalsOperation) IsSplitOperation() bool {
	return m != nil && m.Op == "y" && m.RevealInputIndex == 0
}

// is_seal_operation
func (m *WitnessAtomicalsOperation) Is_seal_operation() bool {
	return m != nil && m.Op == "sl" && m.RevealInputIndex == 0
}

// is_event_operation
func (m *WitnessAtomicalsOperation) Is_event_operation() bool {
	return m != nil && m.Op == "evt" && m.RevealInputIndex == 0
}

// is_immutable
func (m *PayLoad) IsImmutable() bool {
	if m == nil {
		return false
	}
	if m.Args == nil {
		return false
	}
	return m.Args.I
}

func (m *PayLoad) CheckRequest() bool {
	if m.Args == nil {
		return false
	} else {
		request_counter := 0 // # Ensure that only one of the following may be requested || fail
		if m.Args.RequestRealm != "" {
			request_counter += 1
		}
		if m.Args.RequestSubRealm != "" {
			request_counter += 1
		}
		if m.Args.RequestContainer != "" {
			request_counter += 1
		}
		if m.Args.RequestTicker != "" {
			request_counter += 1
		}
		if m.Args.RequestDmitem != "" {
			request_counter += 1
		}
		if request_counter > 1 {
			return false
		}
	}
	return true
}
