package btcsync

import (
	"sync"

	"github.com/btcsuite/btcd/rpcclient"
)

const BlockCacheNum = 3

type BtcSync struct {
	*rpcclient.Client
	blockHeightChannel chan int64
	blockCache         sync.Map
	*TxHeightCache
}

func NewBtcSync(rpcURL, rpcUser, rpcPassword string, startHeight int64) (*BtcSync, error) {
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
	m := &BtcSync{
		Client:             client,
		blockHeightChannel: make(chan int64, 3),
		TxHeightCache:      &TxHeightCache{},
	}
	go m.FetchBlocks()
	return m, nil

}
