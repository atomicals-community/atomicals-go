package svc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atomicals-go/repo/postsql"
	"github.com/atomicals-go/utils"
	"github.com/btcsuite/btcd/btcjson"
)

type PendingAtomicalsAsset struct {
	Mod                    []*postsql.ModInfo
	DeleteFts              []*postsql.UTXOFtInfo
	NewFts                 []*postsql.UTXOFtInfo
	UpdateNfts             []*postsql.UTXONftInfo
	NewUTXOFtInfo          []*postsql.UTXOFtInfo
	UpdateDistributedFt    []*postsql.GlobalDistributedFt
	NewGlobalDistributedFt []*postsql.GlobalDistributedFt
	NewGlobalDirectFt      []*postsql.GlobalDirectFt
	NewUTXONftInfo         []*postsql.UTXONftInfo
}

func (m *PendingAtomicalsAsset) CheckAssetByLocationID(locationID string) string {
	assets := ""
	for _, v := range m.NewFts {
		if v.LocationID == locationID {
			res, err := json.Marshal(v)
			if err != nil {
				continue
			}
			assets += string(res)
		}
	}
	for _, v := range m.UpdateNfts {
		if v.LocationID == locationID {
			res, err := json.Marshal(v)
			if err != nil {
				continue
			}
			assets += string(res)
		}
	}
	for _, v := range m.NewUTXOFtInfo {
		if v.LocationID == locationID {
			res, err := json.Marshal(v)
			if err != nil {
				continue
			}
			assets += string(res)
		}
	}
	for _, v := range m.NewUTXONftInfo {
		if v.LocationID == locationID {
			res, err := json.Marshal(v)
			if err != nil {
				continue
			}
			assets += string(res)
		}
	}
	return assets
}
func newPendingAtomicalsAsset() *PendingAtomicalsAsset {
	a := &PendingAtomicalsAsset{}
	a.Mod = make([]*postsql.ModInfo, 0)
	a.DeleteFts = make([]*postsql.UTXOFtInfo, 0)
	a.NewFts = make([]*postsql.UTXOFtInfo, 0)
	a.UpdateNfts = make([]*postsql.UTXONftInfo, 0)
	a.NewUTXOFtInfo = make([]*postsql.UTXOFtInfo, 0)
	a.UpdateDistributedFt = make([]*postsql.GlobalDistributedFt, 0)
	a.NewGlobalDistributedFt = make([]*postsql.GlobalDistributedFt, 0)
	a.NewGlobalDirectFt = make([]*postsql.GlobalDirectFt, 0)
	a.NewUTXONftInfo = make([]*postsql.UTXONftInfo, 0)
	return a
}

func (a *ServiceContext) SyncMempoolAtomicalsAsset(tx btcjson.TxRawResult, height int64) *PendingAtomicalsAsset {
	pendingAtomicalsAsset := newPendingAtomicalsAsset()

	Mod, DeleteFts, NewFts, UpdateNfts, NewUTXOFtInfo,
		UpdateDistributedFt, NewGlobalDistributedFt, NewGlobalDirectFt, NewUTXONftInfo := a.TraceTx(tx, height)

	pendingAtomicalsAsset.Mod = append(pendingAtomicalsAsset.Mod, Mod)
	pendingAtomicalsAsset.DeleteFts = append(pendingAtomicalsAsset.DeleteFts, DeleteFts...)
	pendingAtomicalsAsset.NewFts = append(pendingAtomicalsAsset.NewFts, NewFts...)
	pendingAtomicalsAsset.UpdateNfts = append(pendingAtomicalsAsset.UpdateNfts, UpdateNfts...)
	pendingAtomicalsAsset.NewUTXOFtInfo = append(pendingAtomicalsAsset.NewUTXOFtInfo, NewUTXOFtInfo)
	pendingAtomicalsAsset.UpdateDistributedFt = append(pendingAtomicalsAsset.UpdateDistributedFt, UpdateDistributedFt)
	pendingAtomicalsAsset.NewGlobalDistributedFt = append(pendingAtomicalsAsset.NewGlobalDistributedFt, NewGlobalDistributedFt)
	pendingAtomicalsAsset.NewGlobalDirectFt = append(pendingAtomicalsAsset.NewGlobalDirectFt, NewGlobalDirectFt)
	pendingAtomicalsAsset.NewUTXONftInfo = append(pendingAtomicalsAsset.NewUTXONftInfo, NewUTXONftInfo)
	return pendingAtomicalsAsset
}

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
		pendingAtomicalsAsset := newPendingAtomicalsAsset()
		for height := int64(location.BlockHeight + 1); height <= maxBlockHeight; height++ {
			block, err := a.GetBlockByHeight(height)
			if err != nil {
				return
			}
			for index := int64(location.TxIndex + 1); index < int64(len(block.Tx)); index++ {
				Mod, DeleteFts, NewFts, UpdateNfts, NewUTXOFtInfo,
					UpdateDistributedFt, NewGlobalDistributedFt, NewGlobalDirectFt, NewUTXONftInfo := a.TraceTx(block.Tx[index], height)
				pendingAtomicalsAsset.Mod = append(pendingAtomicalsAsset.Mod, Mod)
				pendingAtomicalsAsset.DeleteFts = append(pendingAtomicalsAsset.DeleteFts, DeleteFts...)
				pendingAtomicalsAsset.NewFts = append(pendingAtomicalsAsset.NewFts, NewFts...)
				pendingAtomicalsAsset.UpdateNfts = append(pendingAtomicalsAsset.UpdateNfts, UpdateNfts...)
				pendingAtomicalsAsset.NewUTXOFtInfo = append(pendingAtomicalsAsset.NewUTXOFtInfo, NewUTXOFtInfo)
				pendingAtomicalsAsset.UpdateDistributedFt = append(pendingAtomicalsAsset.UpdateDistributedFt, UpdateDistributedFt)
				pendingAtomicalsAsset.NewGlobalDistributedFt = append(pendingAtomicalsAsset.NewGlobalDistributedFt, NewGlobalDistributedFt)
				pendingAtomicalsAsset.NewGlobalDirectFt = append(pendingAtomicalsAsset.NewGlobalDirectFt, NewGlobalDirectFt)
				pendingAtomicalsAsset.NewUTXONftInfo = append(pendingAtomicalsAsset.NewUTXONftInfo, NewUTXONftInfo)
			}
		}
		a.PendingAtomicalsAsset = pendingAtomicalsAsset
		a.SyncHeight = location.BlockHeight
		a.MaxBlockHeight = maxBlockHeight
		time.Sleep(10 * time.Minute)
	}
}
