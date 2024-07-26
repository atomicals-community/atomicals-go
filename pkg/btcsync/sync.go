package btcsync

import (
	"sync"

	"github.com/btcsuite/btcd/rpcclient"
)

const BlockCacheNum = 3

type BtcSync struct {
	*rpcclient.Client
	CurrentHeight      int64
	blockHeightChannel chan int64
	blockCache         sync.Map
	txCache            map[string]int64   // key: txid, value blockheight
	txCacheByHeight    map[int64][]string // key: blockheight, value txid
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
	m := &BtcSync{
		Client:             client,
		blockHeightChannel: make(chan int64, BlockCacheNum),
		txCache:            make(map[string]int64),
		txCacheByHeight:    make(map[int64][]string, 0),
	}
	go m.FetchBlocks()
	return m, nil
}
