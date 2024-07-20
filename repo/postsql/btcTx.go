package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const atomicalsTxTableName = "atomicals_tx"

type AtomicalsTx struct {
	gorm.Model
	BlockHeight int64 `gorm:"index"`
	TxIndex     int64
	TxID        string `gorm:"uniqueindex"`
	Operation   string
	Description string
}

func (*AtomicalsTx) TableName() string {
	return atomicalsTxTableName
}

func (*AtomicalsTx) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(atomicalsTxTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&AtomicalsTx{})
	assert.Nil(nil, err)
}

func (*AtomicalsTx) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(AtomicalsTx{})
	assert.Nil(nil, err)
}
