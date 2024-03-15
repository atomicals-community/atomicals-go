package witness

import (
	"encoding/hex"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/fxamacker/cbor/v2"
)

type WitnessAtomicalsOperation struct {
	Op      string
	Payload *PayLoad
	Script  string

	CommitTxID      string // vin's txID
	CommitVoutIndex int64  // vin's index as vout in last tx
	AtomicalsID     string
	CommitHeight    int64

	RevealLocationTxID      string
	RevealInputIndex        int64
	RevealLocationVoutIndex int64 // always is expect VOUT_EXPECT_OUTPUT_INDEX
	RevealLocationHeight    int64
}

// is_dft_bitwork_rollover_activated
func (m *WitnessAtomicalsOperation) IsDftBitworkRolloverActivated() bool {
	return m.RevealLocationHeight >= common.ATOMICALS_ACTIVATION_HEIGHT_DFT_BITWORK_ROLLOVER
}

// is_within_acceptable_blocks_for_name_reveal
func (m *WitnessAtomicalsOperation) IsWithinAcceptableBlocksForNameReveal() bool {
	return m.CommitHeight >= m.RevealLocationHeight-common.MINT_REALM_CONTAINER_TICKER_COMMIT_REVEAL_DELAY_BLOCKS
}

// is_within_acceptable_blocks_for_general_reveal
func (m *WitnessAtomicalsOperation) IsWithinAcceptableBlocksForGeneralReveal() bool {
	return m.CommitHeight >= m.RevealLocationHeight-common.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS
}

// is_immutable
func (m *WitnessAtomicalsOperation) IsImmutable() bool {
	if m.Payload == nil {
		return false
	}
	if m.Payload.Args == nil {
		return false
	}
	return m.Payload.Args.I
}

func (m *WitnessAtomicalsOperation) IsValidBitwork() (*common.Bitwork, *common.Bitwork, error) {
	if m.Payload == nil {
		return nil, nil, nil
	}
	if m.Payload.Args == nil {
		return nil, nil, nil
	}
	bitworkc := common.ParseBitwork(m.Payload.Args.Bitworkc)
	if bitworkc != nil {
		if !common.IsProofOfWorkPrefixMatch(m.CommitTxID, bitworkc.Prefix, bitworkc.Ext) {
			return nil, nil, errors.ErrInvalidBitWork
		}
	}
	bitworkr := common.ParseBitwork(m.Payload.Args.Bitworkr)
	if bitworkr != nil {
		if !common.IsProofOfWorkPrefixMatch(m.CommitTxID, bitworkr.Prefix, bitworkr.Ext) {
			return nil, nil, errors.ErrInvalidBitWork
		}
	}
	return bitworkc, bitworkr, nil
}

func IsValidMintBitwork(commitTxID, mintBitworkc, mintBitworkr string) (*common.Bitwork, *common.Bitwork, error) {
	bitworkc := common.ParseBitwork(mintBitworkc)
	if bitworkc != nil {
		if !common.IsProofOfWorkPrefixMatch(commitTxID, bitworkc.Prefix, bitworkc.Ext) {
			return nil, nil, errors.ErrInvalidBitWork
		}
	}
	bitworkr := common.ParseBitwork(mintBitworkr)
	if bitworkr != nil {
		if !common.IsProofOfWorkPrefixMatch(commitTxID, bitworkr.Prefix, bitworkr.Ext) {
			return nil, nil, errors.ErrInvalidBitWork
		}
	}
	return bitworkc, bitworkr, nil
}

// is_splat_operation
func (m *WitnessAtomicalsOperation) IsSplatOperation() bool {
	return m != nil && m.Op == "x" && m.CommitVoutIndex == 0
}

// is_split_operation
func (m *WitnessAtomicalsOperation) IsSplitOperation() bool {
	return m != nil && m.Op == "y" && m.CommitVoutIndex == 0
}

// # Parses and detects valid Atomicals protocol operations in a witness script
// # Stops when it finds the first operation in the first input
func ParseWitness(tx btcjson.TxRawResult, height int64) *WitnessAtomicalsOperation {
	for vinIndex, vin := range tx.Vin {
		if !vin.HasWitness() {
			continue
		}
		for _, script := range vin.Witness {
			op, payload, err := ParseOperationAndPayLoad(script)
			if err != nil {
				continue
			}
			return &WitnessAtomicalsOperation{
				Op:                      op,
				Payload:                 payload,
				Script:                  script,
				CommitTxID:              vin.Txid,
				CommitVoutIndex:         int64(vin.Vout),
				AtomicalsID:             common.AtomicalsID(vin.Txid, int64(vin.Vout)),
				RevealLocationTxID:      tx.Txid,
				RevealInputIndex:        int64(vinIndex),
				RevealLocationVoutIndex: common.VOUT_EXPECT_OUTPUT_INDEX,
				RevealLocationHeight:    height,
			}
		}
	}
	return &WitnessAtomicalsOperation{
		RevealLocationTxID:      tx.Txid,
		RevealInputIndex:        -1,
		RevealLocationVoutIndex: common.VOUT_EXPECT_OUTPUT_INDEX,
		RevealLocationHeight:    height,
	}
}

func ParseOperationAndPayLoad(script string) (string, *PayLoad, error) {
	scriptBytes, err := hex.DecodeString(script)
	if err != nil {
		return "", nil, err
	}
	scriptEntryLen := int64(len(scriptBytes))
	if scriptEntryLen < 39 || scriptBytes[0] != 0x20 {
		return "", nil, errors.ErrInvalidWitnessScriptLength
	}
	pkFlag := scriptBytes[0]
	if pkFlag != 0x20 {
		return "", nil, errors.ErrInvalidWitnessScriptPkFlag
	}
	// TODO: the loop below is so confused. procotal should give specific index range
	for index := int64(33); index < scriptEntryLen-6; index++ {
		opFlag := scriptBytes[index]
		if opFlag != OP_IF {
			continue
		}
		if hex.EncodeToString(scriptBytes[index+1:index+6]) != common.ATOMICALS_ENVELOPE_MARKER_BYTES {
			continue
		}
		operation, startIndex := parseAtomicalsOperation(scriptBytes, index+6)
		if operation == "" {
			continue
		}
		payloadBytes, err := parseAtomicalsData(scriptBytes, startIndex)
		if err != nil {
			return "", nil, err
		}
		if payloadBytes == nil {
			continue
		}
		payload := &PayLoad{}
		if err := cbor.Unmarshal(payloadBytes, payload); err != nil {
			return "", nil, err
		}
		if !payload.check() {
			return "", nil, errors.ErrInvalidPayLoad
		}
		// payloadstr, _ := json.Marshal(payload)
		// if operation == "dft" {
		// 	log.Log.Warnf("script:%+v", script)
		// }
		// if payload.Args.RequestContainer != "" {
		// 	log.Log.Warnf("script:%+v", script)
		// 	log.Log.Warnf("payload:%+v", string(payloadstr))
		// }
		// if payload.Args.RequestDmitem != "" {
		// 	log.Log.Warnf("script:%+v", script)
		// 	log.Log.Warnf("payload:%+v", string(payloadstr))
		// }
		// if payload.Args.RequestSubRealm != "" {
		// 	log.Log.Warnf("script:%+v", script)
		// 	log.Log.Warnf("payload:%+v", string(payloadstr))
		// }
		// if payload.Args.RequestRealm != "" {
		// 	log.Log.Warnf("script:%+v", script)
		// 	log.Log.Warnf("payload:%+v", string(payloadstr))
		// }
		return operation, payload, nil
	}
	return "", nil, errors.ErrOptionNotFound
}
