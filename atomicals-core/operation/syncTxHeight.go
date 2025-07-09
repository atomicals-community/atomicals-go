package atomicals

import (
	"github.com/atomicals-go/atomicals-core/witness"
)

func (m *Atomicals) SyncTxHeight() {
	height := m.location.BlockHeight + 2
	for {
		m.syncTxHeight(height)
		height++
	}
}
func (m *Atomicals) syncTxHeight(height int64) {
	block, err := m.GetBlockByHeight(height)
	if err != nil {
		return
	}
	for _, tx := range block.Tx {
		operation := witness.ParseWitness(tx, block.Height)
		if operation.Op == "" {
			continue
		}
		commitHeight, err := m.GetTxHeightByTxID(operation.CommitTxID)
		if err != nil {
			panic(err)
		}
		if _, ok := m.SyncTxHeightMap.Load(operation.CommitTxID); !ok {
			m.SyncTxHeightMap.Store(operation.CommitTxID, commitHeight)
		}
	}
}

func (m *Atomicals) getTxHeight(commitTxID string) int64 {
	h, err := m.GetTxHeightByTxID(commitTxID)
	if err != nil {
		panic(err)
	}
	return h
	// height, ok := m.SyncTxHeightMap.LoadAndDelete(commitTxID)
	// if ok {
	// 	h, _ := height.(int64)
	// 	return h
	// } else {
	// 	var err error
	// 	h, err := m.GetTxHeightByTxID(commitTxID)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	m.SyncTxHeightMap.Delete(commitTxID)
	// 	return h
	// }
}
