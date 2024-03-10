package main

import (
	"time"

	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/btcsync"
	"github.com/atomicals-core/pkg/log"
)

func main() {
	a := atomicals.NewAtomicals(common.ATOMICALS_ACTIVATION_HEIGHT)
	b, err := btcsync.NewBtcSync("rpcURL", "user", "password")
	if err != nil {
		panic(err)
	}
	tracedAllTx(a, b)
	// tracedSpecificTx(a, b, "384b707797e10637b9c3c5f971b3beee0b7cd531ab7c14fda6a320e4b19d4b1f", 0)
}

func tracedAllTx(a *atomicals.Atomicals, b *btcsync.BtcSync) {
	for {
		startTime := time.Now()
		blockInfo, err := b.GetBlockByHeight(a.Height)
		if err != nil {
			log.Log.Warnf("GetBlockByHeight err:%v height:%v", err, a.Height)
			panic(err)
		}
		a.TraceBlock(blockInfo)
		log.Log.Warnf("time.Since(startTime):%v", time.Since(startTime))
	}
}

func tracedSpecificTx(a *atomicals.Atomicals, b *btcsync.BtcSync, txHash string, height int64) {
	tx, err := b.GetTransaction(txHash)
	if err != nil {
		log.Log.Warnf("GetTransaction err:%v", err)
	}
	a.TraceTx(*tx, height)
}
