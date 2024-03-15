package atomicals

import (
	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/pkg/btcsync"
	"github.com/atomicals-core/pkg/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Atomicals struct {
	btcClient *btcsync.BtcSync
	db.DB
}

func NewAtomicalsWithMemory(conf *conf.Config, height int64) *Atomicals {
	b, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	return &Atomicals{
		DB:        db.NewMemoryDB(height, ""),
		btcClient: b,
	}
}

func NewAtomicalsWithSQL(conf *conf.Config, height int64) *Atomicals {
	b, err := btcsync.NewBtcSync(conf.BtcRpcURL, conf.BtcRpcUser, conf.BtcRpcPassword)
	if err != nil {
		panic(err)
	}
	DB, err := gorm.Open(postgres.Open(conf.SqlDNS), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	return &Atomicals{
		DB:        db.NewSqlDB(DB, height, ""),
		btcClient: b,
	}
}
