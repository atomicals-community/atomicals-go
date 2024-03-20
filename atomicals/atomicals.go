package atomicals

import (
	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/pkg/btcsync"
	"github.com/atomicals-core/pkg/conf"
)

type Atomicals struct {
	*btcsync.BtcSync
	db.DB
}

func NewAtomicalsWithMemory(conf *conf.Config) *Atomicals {
	b, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	return &Atomicals{
		// DB:        db.NewMemoryDB(),
		BtcSync: b,
	}
}

func NewAtomicalsWithSQL(conf *conf.Config) *Atomicals {
	b, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	d := db.NewSqlDB(conf.SqlDNS, b)
	return &Atomicals{
		DB:      d,
		BtcSync: b,
	}
}
