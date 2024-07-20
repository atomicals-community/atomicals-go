package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const LocationTableName = "location"

type Location struct {
	gorm.Model
	Name        string
	BlockHeight int64
	TxIndex     int64
	Txid        string
}

func (*Location) TableName() string {
	return LocationTableName
}

func (*Location) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(LocationTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&Location{})
	assert.Nil(nil, err)
	dbTx := db.Save(&Location{Name: "atomicals", BlockHeight: 812480, TxIndex: -1})
	assert.Nil(nil, dbTx.Error)
}

func (*Location) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(Location{})
	assert.Nil(nil, err)
}
