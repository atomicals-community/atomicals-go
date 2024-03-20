package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const UserFtInfoTableName = "atomicals_user_ft_utxo"

type UserFtInfo struct {
	gorm.Model
	UserPk      string
	AtomicalsID string `gorm:"index"` // (txID,VOUT_EXPECT_OUTPUT_INDEX) init when be minted
	LocationID  string `gorm:"index"` // (txID,voutIndex)updated after being transfered
	// Bitworkc    *common.Bitwork
	// Bitworkr    *common.Bitwork

	// DistributedFt
	MintTicker string `gorm:"index"`
	Nonce      int64
	Time       int64
	// MintBitworkVec  *common.Bitwork
	// MintBitworkcInc *common.Bitwork
	// MintBitworkrInc *common.Bitwork
	Amount int64

	// DirectFt
	Type          string
	Subtype       string
	RequestTicker string
	MaxSupply     int64
	MintAmount    int64
	MintHeight    int64
	MaxMints      int64
	// Meta          *witness.Meta
}

func (*UserFtInfo) TableName() string {
	return UserFtInfoTableName
}

func (*UserFtInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(UserFtInfoTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&UserFtInfo{})
	assert.Nil(nil, err)
}

func (*UserFtInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(UserFtInfo{})
	assert.Nil(nil, err)
}
