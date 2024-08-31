package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const LocationTableName = "location"
const LocationKey = "atomicals"

type Location struct {
	gorm.Model
	Key         string `gorm:"uniqueindex"`
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
}

func (*Location) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(Location{})
	assert.Nil(nil, err)
}
