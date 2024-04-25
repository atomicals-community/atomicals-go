package api

import "github.com/atomicals-core/atomicals/DB/postsql"

type ReqAssetByAtomicalsID struct {
	AtomicalsID string `json:"atomicals_id"`
}

type RespNftByAtomicalsID struct {
	Assets []*postsql.UTXONftInfo `json:"assets"`
}

type RespFtByAtomicalsID struct {
	Assets []*postsql.UTXOFtInfo `json:"assets"`
}
