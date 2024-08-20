package atomicals

import (
	"github.com/atomicals-go/atomicals-indexer/atomicals-core/witness"
	"github.com/atomicals-go/repo/postsql"
	"github.com/btcsuite/btcd/btcjson"
)

func (m *Atomicals) operationDat(operation *witness.WitnessAtomicalsOperation, tx btcjson.TxRawResult) *postsql.DatInfo {
	return &postsql.DatInfo{
		Height:      operation.RevealLocationHeight,
		AtomicalsID: operation.AtomicalsID,
		LocationID:  operation.LocationID,
		Dat:         operation.PayloadStr,
	}
}
