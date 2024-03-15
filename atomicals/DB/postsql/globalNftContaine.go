package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const globalNftContainerTableName = "atomicals_nft_utxo"

type GlobalNftContainer struct {
	gorm.Model
	ContainerName string `gorm:"uniqueindex"`
	ItemID        string `gorm:"index"`
}

func (*GlobalNftContainer) TableName() string {
	return globalNftContainerTableName
}

func (*GlobalNftContainer) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(globalNftContainerTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&GlobalNftContainer{})
	assert.Nil(nil, err)
}

func (*GlobalNftContainer) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(GlobalNftContainer{})
	assert.Nil(nil, err)
}
