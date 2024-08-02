package atomicals

import (
	"github.com/atomicals-go/pkg/btcsync"
	"github.com/atomicals-go/pkg/conf"
	"github.com/atomicals-go/pkg/log"
	"github.com/atomicals-go/repo"
	"github.com/atomicals-go/repo/postsql"
)

type Atomicals struct {
	*btcsync.BtcSync
	repo.DB
	location       *postsql.Location
	maxBlockHeight int64
}

func NewAtomicalsWithSQL(conf *conf.Config) *Atomicals {
	db := repo.NewSqlDB(conf.SqlDNS)
	location, err := db.Location()
	if err != nil {
		log.Log.Panicf("Location err:%v", err)
	}
	btcsync, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	maxBlockHeight, err := btcsync.GetBlockCount()
	if err != nil {
		log.Log.Panicf("GetBlockCount err:%v", err)
	}
	return &Atomicals{
		DB:             db,
		BtcSync:        btcsync,
		location:       location,
		maxBlockHeight: maxBlockHeight,
	}
}
