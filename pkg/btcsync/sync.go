package btcsync

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

type BtcSync struct {
	*rpcclient.Client
}

func NewBtcSync(rpcURL, rpcUser, rpcPassword string) (*BtcSync, error) {
	connCfg := &rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         rpcURL,
		// Endpoint:     rpcURL,
		User: rpcUser,
		Pass: rpcPassword,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}
	return &BtcSync{
		client,
	}, nil
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

func (m *BtcSync) GetBlockByHeight(blockHeight int64) (*btcjson.GetBlockVerboseTxResult, error) {
	blockHash, err := m.GetBlockHash(blockHeight)
	if err != nil {
		return nil, err
	}
	block, err := m.GetBlockVerboseTx(blockHash)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (m *BtcSync) GetCommitHeight(txID string) (int64, error) {
	txHash, err := TxID2txHash(txID)
	if err != nil {
		return -1, err
	}
	t, err := m.GetTransaction(txHash)
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

func TxID2txHash(txID string) (string, error) {
	txidBytes, err := hex.DecodeString(txID)
	if err != nil {
		return "", err
	}
	txhash := hex.EncodeToString(txidBytes)
	return txhash, nil
}
