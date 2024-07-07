package svc

import (
	"fmt"
	"time"

	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

func (a *ServiceContext) SyncPendingAtomicalsAsset() {
	location, err := a.Location()
	if err != nil {
		panic(err)
	}
	maxBlockHeight, err := a.GetBlockCount()
	if err != nil {
		panic(err)
	}
	if location.BlockHeight < maxBlockHeight-utils.SafeBlockHeightInterupt {
		panic(fmt.Sprintf("waiting for atomicals-core sync to:%v, current height:%v", maxBlockHeight-utils.SafeBlockHeightInterupt, location.BlockHeight))
	}
	if a.SyncHeight != location.BlockHeight {
		pendingAtomicalsAssetMap := make(map[string]*PendingAtomicalsAsset, 0)
		for height := int64(location.BlockHeight + 1); height <= maxBlockHeight; height++ {
			block, err := a.GetBlockByHeight(height)
			if err != nil {
				return
			}
			for index := int64(location.TxIndex + 1); index < int64(len(block.Tx)); index++ {
				tx := block.Tx[index]
				Mod, DeleteFts, NewFts, UpdateNfts, NewUTXOFtInfo,
					UpdateDistributedFt, NewGlobalDistributedFt, NewGlobalDirectFt, NewUTXONftInfo := a.TraceTx(tx, height)
				pendingAtomicalsAsset := newPendingAtomicalsAsset()
				pendingAtomicalsAsset.Mod = Mod
				pendingAtomicalsAsset.DeleteFts = DeleteFts
				pendingAtomicalsAsset.NewFts = NewFts
				pendingAtomicalsAsset.UpdateNfts = UpdateNfts
				pendingAtomicalsAsset.NewUTXOFtInfo = NewUTXOFtInfo
				pendingAtomicalsAsset.UpdateDistributedFt = UpdateDistributedFt
				pendingAtomicalsAsset.NewGlobalDistributedFt = NewGlobalDistributedFt
				pendingAtomicalsAsset.NewGlobalDirectFt = NewGlobalDirectFt
				pendingAtomicalsAsset.NewUTXONftInfo = NewUTXONftInfo
				pendingAtomicalsAssetMap[tx.Txid] = pendingAtomicalsAsset

			}
		}
		a.PendingAtomicalsAssetMap = pendingAtomicalsAssetMap
		a.SyncHeight = location.BlockHeight
		a.MaxBlockHeight = maxBlockHeight
		time.Sleep(10 * time.Minute)
	}
}

func (a *ServiceContext) SyncMempoolAtomicalsAsset(tx btcjson.TxRawResult, height int64) *PendingAtomicalsAsset {
	pendingAtomicalsAsset := newPendingAtomicalsAsset()

	Mod, DeleteFts, NewFts, UpdateNfts, NewUTXOFtInfo,
		UpdateDistributedFt, NewGlobalDistributedFt, NewGlobalDirectFt, NewUTXONftInfo := a.TraceTx(tx, height)
	pendingAtomicalsAsset.Mod = Mod
	pendingAtomicalsAsset.DeleteFts = DeleteFts
	pendingAtomicalsAsset.NewFts = NewFts
	pendingAtomicalsAsset.UpdateNfts = UpdateNfts
	pendingAtomicalsAsset.NewUTXOFtInfo = NewUTXOFtInfo
	pendingAtomicalsAsset.UpdateDistributedFt = UpdateDistributedFt
	pendingAtomicalsAsset.NewGlobalDistributedFt = NewGlobalDistributedFt
	pendingAtomicalsAsset.NewGlobalDirectFt = NewGlobalDirectFt
	pendingAtomicalsAsset.NewUTXONftInfo = NewUTXONftInfo
	return pendingAtomicalsAsset
}
