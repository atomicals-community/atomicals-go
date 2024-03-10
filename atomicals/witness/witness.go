package witness

import (
	"encoding/hex"
	"log"

	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/errors"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/fxamacker/cbor/v2"
)

type WitnessAtomicalsOperation struct {
	Op       string
	Payload  *PayLoad
	TxID     string
	VinIndex int64
	Height   int64
}

// is_splat_operation
func (m *WitnessAtomicalsOperation) IsSplatOperation() bool {
	return m != nil && m.Op == "x" && m.VinIndex == 0
}

// is_split_operation
func (m *WitnessAtomicalsOperation) IsSplitOperation() bool {
	return m != nil && m.Op == "y" && m.VinIndex == 0
}

// # Parses and detects valid Atomicals protocol operations in a witness script
// # Stops when it finds the first operation in the first input
func ParseOperationAndPayLoad(tx btcjson.TxRawResult, height int64) *WitnessAtomicalsOperation {
	for vinIndex, vin := range tx.Vin {
		if !vin.HasWitness() {
			continue
		}
		for _, script := range vin.Witness {
			op, payload, err := parseOperationAndPayLoad(script)
			if err != nil {
				log.Printf("parseOperationAndPayLoad err:%+v", err)
				continue
			}
			return &WitnessAtomicalsOperation{
				Op:       op,
				Payload:  payload,
				TxID:     tx.Txid,
				VinIndex: int64(vinIndex),
				Height:   height,
			}
		}
	}
	return &WitnessAtomicalsOperation{
		TxID:   tx.Txid,
		Height: height,
	}
}

func parseOperationAndPayLoad(script string) (string, *PayLoad, error) {
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
		if opFlag != common.OP_IF {
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
		return operation, payload, nil
	}
	return "", nil, errors.ErrOptionNotFound
}
