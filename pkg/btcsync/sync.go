package btcsync

import (
	"sync"

	"github.com/btcsuite/btcd/rpcclient"
)

type BtcSync struct {
	*rpcclient.Client
	blockCacheNum      int
	blockHeightChannel chan int64
	blockCache         sync.Map

	*TxHeightCache
}

func NewBtcSync(rpcURL, rpcUser, rpcPassword string) (*BtcSync, error) {
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         rpcURL,
		User:         rpcUser,
		Pass:         rpcPassword,
	}, nil)
	if err != nil {
		return nil, err
	}
	b := &BtcSync{
		Client:             client,
		blockCacheNum:      5,
		blockHeightChannel: make(chan int64, 5),
		TxHeightCache:      newTxHeightCache(),
	}
	go b.FetchBlocks()

	return b, nil

}
