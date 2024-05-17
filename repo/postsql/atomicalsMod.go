package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const ModTableName = "atomicals_mod"

type ModInfo struct {
	gorm.Model
	AtomicalsID string `gorm:"index"` // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	LocationID  string `gorm:"index"` // txID_voutIndex updated after being transfered
	Mod         string
}

func (*ModInfo) TableName() string {
	return ModTableName
}

func (*ModInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(ModTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&ModInfo{})
	assert.Nil(nil, err)
}

func (*ModInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(ModInfo{})
	assert.Nil(nil, err)
}
