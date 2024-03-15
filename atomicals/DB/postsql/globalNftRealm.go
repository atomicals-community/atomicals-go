package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const GlobalNftRealmTableName = "atomicals_nft_utxo"

type GlobalNftRealm struct {
	gorm.Model
	RealmName    string `gorm:"uniqueindex"`
	SubRealmName string `gorm:"index"` // when realm was init, default SubRealmName is ""
}

func (*GlobalNftRealm) TableName() string {
	return GlobalNftRealmTableName
}

func (*GlobalNftRealm) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(GlobalNftRealmTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&GlobalNftRealm{})
	assert.Nil(nil, err)
}

func (*GlobalNftRealm) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(GlobalNftRealm{})
	assert.Nil(nil, err)
}
