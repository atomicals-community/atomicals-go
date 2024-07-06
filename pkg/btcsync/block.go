package btcsync

import (
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func (m *BtcSync) GetBlockByHeightSync(blockHeight int64) *btcjson.GetBlockVerboseTxResult {
	for height := blockHeight; height < blockHeight+int64(BlockCacheNum); height++ {
		if m.CurrentHeight < height {
			m.blockHeightChannel <- height
			m.CurrentHeight = height
		}
	}
	return m.BlockByHeight(blockHeight)
}

func (m *BtcSync) BlockByHeight(blockHeight int64) *btcjson.GetBlockVerboseTxResult {
	var b *btcjson.GetBlockVerboseTxResult
	for {
		block, ok := m.blockCache.Load(blockHeight)
		if ok {
			b, _ = block.(*btcjson.GetBlockVerboseTxResult)
			m.blockCache.Delete(blockHeight - BlockCacheNum)
			break
		}
		time.Sleep(1 * time.Second)
	}
	return b
}

func (m *BtcSync) FetchBlocks() error {
	for height := range m.blockHeightChannel {
		// set block cache
		block, err := m.GetBlockByHeight(height)
		if err != nil {
			continue
		}
		m.blockCache.Store(height, block)
	}
	return nil
}

func (m *BtcSync) GetBlockByHeight(height int64) (*btcjson.GetBlockVerboseTxResult, error) {
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

func (m *BtcSync) GetBlockCount() (int64, error) {
	return m.Client.GetBlockCount()
}

func (m *BtcSync) GetBlockHash(blockHeight int64) (*chainhash.Hash, error) {
	return m.Client.GetBlockHash(blockHeight)
}
