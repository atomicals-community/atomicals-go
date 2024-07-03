package btcsync

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func (m *BtcSync) GetTxHeightByTxID(txID string) (int64, error) {
	t, err := m.GetTransaction(txID)
	if err != nil {
		return -1, err
	}
	blockHash, err := chainhash.NewHashFromStr(t.BlockHash)
	if err != nil {
		return -1, err
	}
	blockInfo, err := m.GetBlockVerboseTx(blockHash)
	if err != nil {
		return -1, err
	}
	return blockInfo.Height, nil
}

func (m *BtcSync) GetTransaction(txHash string) (*btcjson.TxRawResult, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, err
	}
	rawTx, err := m.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, err
	}
	return rawTx, nil
}
