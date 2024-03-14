package atomicals

import (
	db "github.com/atomicals-core/atomicals/DB"
	"github.com/atomicals-core/pkg/btcsync"
)

type Atomicals struct {
	btcClient *btcsync.BtcSync
	db.DB
}

func NewAtomicals(btcClient *btcsync.BtcSync, height int64) *Atomicals {
	return &Atomicals{
		DB:        db.NewMemoryDB(height, ""),
		btcClient: btcClient,
	}
}
