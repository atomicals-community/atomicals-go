package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const btcTxTableName = "btc_tx"

type BtcTx struct {
	gorm.Model
	TxID        string `gorm:"uniqueindex"`
	BlockHeight int64
}

func (*BtcTx) TableName() string {
	return btcTxTableName
}

func (*BtcTx) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(btcTxTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&BtcTx{})
	assert.Nil(nil, err)
}

func (*BtcTx) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(BtcTx{})
	assert.Nil(nil, err)
}
