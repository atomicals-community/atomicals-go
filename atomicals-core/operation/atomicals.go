package atomicals

import (
	"time"

	"github.com/atomicals-go/pkg/btcsync"
	"github.com/atomicals-go/pkg/conf"
	"github.com/atomicals-go/repo"
)

type Atomicals struct {
	*btcsync.BtcSync
	repo.DB
	startTime time.Time
}

func NewAtomicalsWithSQL(conf *conf.Config) *Atomicals {
	btcsync, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	db := repo.NewSqlDB(conf.SqlDNS)
	return &Atomicals{
		DB:        db,
		BtcSync:   btcsync,
		startTime: time.Now(),
	}
}
