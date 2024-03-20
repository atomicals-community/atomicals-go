package btcsync

import (
	"errors"
	"time"

	"github.com/btcsuite/btcd/btcjson"
)

func (m *BtcSync) GetBlockByHeight(blockHeight int64) (*btcjson.GetBlockVerboseTxResult, error) {
	for height := int64(blockHeight); height < blockHeight+int64(m.blockCacheNum); height++ {
		m.blockHeightChannel <- height
	}
	block, ok := m.blockCache.LoadAndDelete(blockHeight)
	if ok {
		b, _ := block.(*btcjson.GetBlockVerboseTxResult)
		return b, nil
	}
	time.Sleep(2 * time.Second)
	block, ok = m.blockCache.LoadAndDelete(blockHeight)
	if !ok {
		return nil, errors.ErrUnsupported
	}
	b, _ := block.(*btcjson.GetBlockVerboseTxResult)
	return b, nil
}

func (m *BtcSync) FetchBlocks() error {
	for height := range m.blockHeightChannel {
		if _, ok := m.blockCache.Load(height); ok {
			continue
		}
		blockHash, err := m.GetBlockHash(height)
		if err != nil {
			continue
		}
		block, err := m.GetBlockVerboseTx(blockHash)
		if err != nil {
			continue
		}
		for _, tx := range block.Tx {
			m.SetTxHeightCache(tx.Txid, block.Height)
		}
		m.blockCache.Store(height, block)
	}
	return nil
}
