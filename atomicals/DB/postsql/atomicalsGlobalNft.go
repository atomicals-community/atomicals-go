package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const globalNftTableName = "atomicals_ntf"

type GlobalNft struct {
	gorm.Model
	NftType int64  `gorm:"index"`
	Name    string `gorm:"index"` // realmName containerName tickerName
	SubName string `gorm:"index"` // subRealm  itemID
}

func (*GlobalNft) TableName() string {
	return globalNftTableName
}

func (*GlobalNft) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(globalNftTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&GlobalNft{})
	assert.Nil(nil, err)
}

func (*GlobalNft) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(GlobalNft{})
	assert.Nil(nil, err)
}
