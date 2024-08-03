package svc

import (
	"fmt"
	"time"

	"github.com/atomicals-go/repo"
	"github.com/atomicals-go/utils"
)

func (a *ServiceContext) SyncPendingAtomicalsAsset() {
	location, err := a.Location()
	if err != nil {
		return
	}
	if a.CurrentHeight == location.BlockHeight {
		time.Sleep(10 * time.Minute)
		return
	}
	maxBlockHeight, err := a.GetBlockCount()
	if err != nil {
		return
	}
	if location.BlockHeight < maxBlockHeight-utils.SafeBlockHeightInterupt {
		panic(fmt.Sprintf("waiting for atomicals-core sync to:%v, current height:%v", maxBlockHeight-utils.SafeBlockHeightInterupt, location.BlockHeight))
	}
	pendingAtomicalsAssetMap := make(map[string]*repo.AtomicaslData, 0)
	for height := int64(location.BlockHeight + 1); height <= maxBlockHeight; height++ {
		block, err := a.GetBlockByHeight(height)
		if err != nil {
			return
		}
		for _, tx := range block.Tx {
			data := a.TraceTx(tx, height)
			pendingAtomicalsAssetMap[tx.Txid] = data
		}
	}
	a.PendingAtomicalsAssetMap = pendingAtomicalsAssetMap
	a.CurrentHeight = location.BlockHeight
	a.MaxBlockHeight = maxBlockHeight
}
