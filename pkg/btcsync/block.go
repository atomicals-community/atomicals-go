package btcsync

import (
	"time"

	"github.com/btcsuite/btcd/btcjson"
)

func (m *BtcSync) GetBlockByHeight(blockHeight int64) (*btcjson.GetBlockVerboseTxResult, error) {
	for height := blockHeight; height < blockHeight+int64(BlockCacheNum); height++ {
		if m.CurrentHeight < height {
			m.blockHeightChannel <- height
			m.CurrentHeight = height
		}
	}
	var b *btcjson.GetBlockVerboseTxResult
	for {
		block, ok := m.blockCache.Load(blockHeight)
		if ok {
			b, _ = block.(*btcjson.GetBlockVerboseTxResult)
			break
		}
		m.blockCache.Load(blockHeight - BlockCacheNum)
		time.Sleep(1 * time.Second)
	}
	return b, nil
}

func (m *BtcSync) FetchBlocks() error {
	for height := range m.blockHeightChannel {
		// set block cache
		block, err := m.getBlockByHeight(height)
		if err != nil {
			continue
		}
		m.blockCache.Store(height, block)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (m *BtcSync) getBlockByHeight(height int64) (*btcjson.GetBlockVerboseTxResult, error) {
	blockHash, err := m.GetBlockHash(height)
	if err != nil {
		return nil, err
	}
	block, err := m.GetBlockVerboseTx(blockHash)
	if err != nil {
		return nil, err
	}
	return block, nil
}
