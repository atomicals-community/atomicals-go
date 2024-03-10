package main

import (
	"log"
	"time"

	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/btcsync"
)

func main() {
	a := atomicals.NewAtomicals(common.ATOMICALS_ACTIVATION_HEIGHT)
	b, err := btcsync.NewBtcSync("rpcURL", "user", "password")
	if err != nil {
		panic(err)
	}
	tracedAllTx(a, b)

	// tracedSpecificTx(a, b, "", 0)
}

func tracedAllTx(a *atomicals.Atomicals, b *btcsync.BtcSync) {
	for {
		startTime := time.Now()
		blockInfo, err := b.GetBlockByHeight(a.Height)
		if err != nil {
			log.Printf("GetBlockByHeight err:%v", err)
			panic(err)
		}
		a.TraceBlock(blockInfo)
		log.Printf("time.Since(startTime):%v", time.Since(startTime))
	}
}

func tracedSpecificTx(a *atomicals.Atomicals, b *btcsync.BtcSync, txHash string, height int64) {
	tx, err := b.GetTransaction(txHash)
	if err != nil {
		log.Printf("GetTransaction err:%v", err)
	}
	a.TraceTx(*tx, height)
}
