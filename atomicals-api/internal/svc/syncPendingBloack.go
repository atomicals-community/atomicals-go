package svc

import (
	"encoding/json"

	"github.com/atomicals-go/repo/postsql"
)

type PendingAtomicalsAsset struct {
	Operation              string
	Mod                    *postsql.ModInfo
	DeleteFts              []*postsql.UTXOFtInfo
	NewFts                 []*postsql.UTXOFtInfo
	UpdateNfts             []*postsql.UTXONftInfo
	NewUTXOFtInfo          *postsql.UTXOFtInfo
	UpdateDistributedFt    *postsql.GlobalDistributedFt
	NewGlobalDistributedFt *postsql.GlobalDistributedFt
	NewGlobalDirectFt      *postsql.GlobalDirectFt
	NewUTXONftInfo         *postsql.UTXONftInfo
}

func (m *PendingAtomicalsAsset) CheckAsset() string {
	assets := ""
	for _, v := range m.NewFts {
		res, err := json.Marshal(v)
		if err != nil {
			continue
		}
		assets += string(res)
		m.Operation = "transfer_ft"
	}
	for _, v := range m.UpdateNfts {
		res, err := json.Marshal(v)
		if err != nil {
			continue
		}
		assets += string(res)
		m.Operation = "transfer_nft"
	}
	if m.NewUTXOFtInfo != nil {
		res, _ := json.Marshal(m.NewUTXOFtInfo)
		assets += string(res)
		m.Operation = "ft"
	}
	if m.NewUTXONftInfo != nil {
		res, _ := json.Marshal(m.NewUTXONftInfo)
		assets += string(res)
		m.Operation = "nft"
	}
	return assets
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
	if m.NewUTXOFtInfo != nil {
		res, _ := json.Marshal(m.NewUTXOFtInfo)
		assets += string(res)
	}
	if m.NewUTXONftInfo != nil {
		res, _ := json.Marshal(m.NewUTXONftInfo)
		assets += string(res)
	}
	return assets
}
func newPendingAtomicalsAsset() *PendingAtomicalsAsset {
	a := &PendingAtomicalsAsset{}
	a.DeleteFts = make([]*postsql.UTXOFtInfo, 0)
	a.NewFts = make([]*postsql.UTXOFtInfo, 0)
	a.UpdateNfts = make([]*postsql.UTXONftInfo, 0)
	return a
}
