package btcsync

import (
	"sync"
)

type TxHeightCache struct {
	TxCache         sync.Map // key:txid, value: blockHeight; used for GetCommitHeight
	TxCacheByHeight sync.Map // key: blockHeight, value: map[string]bool(key: txid); used for delete useless txCache(blockHeight<currentHeight-common.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS)
}

func newTxHeightCache() *TxHeightCache {
	return &TxHeightCache{}
}

func (m *BtcSync) GetCommitHeight(txID string) int64 {
	res, ok := m.TxCache.LoadAndDelete(txID)
	// if exist, return height directly
	if ok {
		height, _ := res.(int64)
		return height
	}
	height, err := m.GetTxHeightByTxID(txID)
	if err != nil {
		panic(err)
	}
	return height
}

func (m *BtcSync) SetTxHeightCache(txID string, height int64) {
	m.TxCache.Store(txID, height)
	res, ok := m.TxCacheByHeight.Load(height)
	if ok {
		cache, _ := res.(map[string]bool)
		cache[txID] = true
		m.TxCacheByHeight.Store(height, cache)
	}
	m.TxCacheByHeight.Store(height, map[string]bool{txID: true})
}

// deleteHeight := currentHeight - common.MINT_GENERAL_COMMIT_REVEAL_DELAY_BLOCKS-1
func (m *BtcSync) DeleteUselessTxCache(deleteHeight int64) {
	res, ok := m.TxCacheByHeight.LoadAndDelete(deleteHeight)
	if !ok {
		return
	}
	cache, _ := res.(map[string]bool)
	for txID, _ := range cache {
		m.TxCache.Delete(txID)
	}
}
