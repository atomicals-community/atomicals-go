package witness

import (
	"encoding/hex"

	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/utils"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/fxamacker/cbor/v2"
)

type WitnessAtomicalsOperation struct {
	Script  string
	Op      string
	Payload *PayLoad

	CommitTxID      string // vin's txID
	CommitVoutIndex int64  // vin's index as vout in last tx
	CommitHeight    int64

	AtomicalsID          string
	LocationID           string
	RevealLocationTxID   string
	RevealInputIndex     int64
	RevealLocationHeight int64
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
			// if op == "mod" {
			// 	pythonparse.ParseAtomicalsOperation(script)
			// 	log.Log.Infof("AtomicalsID:%v", utils.AtomicalsID(vin.Txid, int64(vin.Vout)))
			// 	log.Log.Infof("LocationID:%v", utils.AtomicalsID(tx.Txid, int64(vinIndex)))
			// } else if op == "evt" {
			// 	pythonparse.ParseAtomicalsOperation(script)
			// } else if op == "dat" {
			// 	pythonparse.ParseAtomicalsOperation(script)
			// } else if op == "sl" {
			// 	pythonparse.ParseAtomicalsOperation(script)
			// } else if op == "nft" {
			// 	pythonparse.ParseAtomicalsOperation(script)
			// 	log.Log.Infof("AtomicalsID:%v", utils.AtomicalsID(vin.Txid, int64(vin.Vout)))
			// 	log.Log.Infof("LocationID:%v", utils.AtomicalsID(tx.Txid, int64(vinIndex)))
			// 	panic("#")
			// }
			return &WitnessAtomicalsOperation{
				Op:              op,
				Payload:         payload,
				Script:          script,
				CommitTxID:      vin.Txid,
				CommitVoutIndex: int64(vin.Vout),
				CommitHeight:    -1,

				AtomicalsID:          utils.AtomicalsID(vin.Txid, int64(vin.Vout)),
				LocationID:           utils.AtomicalsID(tx.Txid, int64(vinIndex)),
				RevealLocationTxID:   tx.Txid,
				RevealInputIndex:     int64(vinIndex),
				RevealLocationHeight: height,
			}
		}
		break
	}
	return &WitnessAtomicalsOperation{
		RevealLocationTxID:   tx.Txid,
		RevealInputIndex:     -1,
		RevealLocationHeight: height,
	}
}

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

func ParseMintBitwork(commitTxID, mintBitworkc, mintBitworkr string) (*utils.Bitwork, *utils.Bitwork, error) {
	bitworkc := utils.ParseBitwork(mintBitworkc)
	bitworkr := utils.ParseBitwork(mintBitworkr)
	return bitworkc, bitworkr, nil
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
		if hex.EncodeToString(scriptBytes[index+1:index+6]) != utils.ATOMICALS_ENVELOPE_MARKER_BYTES {
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
		// if operation == "nft" {
		// 	if payload.Args.RequestDmitem != "" {
		// 		pythonparse.ParseAtomicalsOperation(script)
		// 	}
		// }
		return operation, payload, nil
	}
	return "", nil, errors.ErrOptionNotFound
}
