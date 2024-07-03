package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const globalDirectFtTableName = "atomicals_global_direct_ft"

type GlobalDirectFt struct {
	gorm.Model
	UserPk      string
	AtomicalsID string `gorm:"index"` // (txID,VOUT_EXPECT_OUTPUT_INDEX) init when be minted
	LocationID  string `gorm:"index"` // (txID,voutIndex)updated after being transfered
	Bitworkc    string
	Bitworkr    string

	Type       string
	Subtype    string
	TickerName string
	MaxSupply  int64
	MintAmount int64
	MintHeight int64
	MaxMints   int64
}

func (*GlobalDirectFt) TableName() string {
	return globalDirectFtTableName
}

func (*GlobalDirectFt) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(globalDirectFtTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&GlobalDirectFt{})
	assert.Nil(nil, err)
}

func (*GlobalDirectFt) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(GlobalDirectFt{})
	assert.Nil(nil, err)
}
