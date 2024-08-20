package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const DatTableName = "atomicals_dat"

type DatInfo struct {
	gorm.Model
	Height      int64
	AtomicalsID string `gorm:"index"` // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	LocationID  string `gorm:"index"` // txID_voutIndex updated after being transfered
	Dat         string
}

func (*DatInfo) TableName() string {
	return DatTableName
}

func (*DatInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(DatTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&DatInfo{})
	assert.Nil(nil, err)
}

func (*DatInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(DatInfo{})
	assert.Nil(nil, err)
}
