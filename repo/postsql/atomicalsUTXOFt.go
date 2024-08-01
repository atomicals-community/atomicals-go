package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const UserFtInfoTableName = "atomicals_utxo_ft"

type FtType string

const (
	TypeDistributedFt FtType = "distributedFt"
	TypeDirectFt      FtType = "directFt"
)

type UTXOFtInfo struct {
	gorm.Model
	UserPk      string
	AtomicalsID string `gorm:"index"` // (txID,VOUT_EXPECT_OUTPUT_INDEX) init when be minted
	LocationID  string `gorm:"index"` // (txID,voutIndex)updated after being transfered
	Bitworkc    string
	Bitworkr    string

	// DistributedFt
	MintTicker      string `gorm:"index"`
	Time            int64
	MintBitworkVec  string
	MintBitworkcInc string
	MintBitworkrInc string
	Amount          int64

	// DirectFt
	Type       string
	Subtype    string
	TickerName string
	MaxSupply  int64
	MintAmount int64
	MintHeight int64
	MaxMints   int64
}

func (*UTXOFtInfo) TableName() string {
	return UserFtInfoTableName
}

func (*UTXOFtInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(UserFtInfoTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&UTXOFtInfo{})
	assert.Nil(nil, err)
}

func (*UTXOFtInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(UTXOFtInfo{})
	assert.Nil(nil, err)
}
