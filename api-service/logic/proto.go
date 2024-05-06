package logic

import "github.com/atomicals-go/atomicals-core/repo/postsql"

type ReqAssetByAtomicalsID struct {
	AtomicalsID string `json:"atomicals_id"`
}

type RespNftByAtomicalsID struct {
	Assets []*postsql.UTXONftInfo `json:"assets"`
}

type RespFtByAtomicalsID struct {
	Assets []*postsql.UTXOFtInfo `json:"assets"`
}
