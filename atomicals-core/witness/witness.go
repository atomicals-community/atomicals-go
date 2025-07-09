package witness

import (
	pythonparse "github.com/atomicals-go/atomicals-core/witness/python-parse"
	"github.com/atomicals-go/utils"

	"github.com/btcsuite/btcd/btcjson"
)

type WitnessAtomicalsOperation struct {
	Script  string
	Op      string
	Payload *PayLoad

	CommitTxID      string // vin's txID
	CommitVoutIndex int64  // vin's index as vout in last tx
	CommitHeight    int64
	CommitTxIndex   int64

	RevealLocationTxID   string
	RevealInputIndex     int64
	RevealLocationHeight int64

	AtomicalsID string
	LocationID  string
}

// # Parses and detects valid Atomicals protocol operations in a witness script
// # Stops when it finds the first operation in the first input
func ParseWitness(tx btcjson.TxRawResult, height int64) *WitnessAtomicalsOperation {
	for vinIndex, vin := range tx.Vin {
		if !vin.HasWitness() {
			continue
		}
		for _, script := range vin.Witness {
			op, payload, err := parseOperationAndPayLoad(script, height)
			if err != nil {
				continue
			}
			if op != "" {
				isValid, err := pythonparse.ParseAtomicalsOperation(script, height)
				if err != nil {
					continue
				}
				if !isValid {
					continue
				}
			}
			return &WitnessAtomicalsOperation{
				Op:              op,
				Payload:         payload,
				Script:          script,
				CommitTxID:      vin.Txid,
				CommitVoutIndex: int64(vin.Vout),
				CommitHeight:    -1,
				CommitTxIndex:   -1,

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
