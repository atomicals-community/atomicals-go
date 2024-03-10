package btcsync

import (
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
