package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const globalDistributedFtTableName = "atomicals_global_distributed_ft"

type GlobalDistributedFt struct {
	gorm.Model
	AtomicalsID    string
	TickerName     string `gorm:"uniqueindex"`
	Type           string
	Subtype        string
	MintMode       string // emu: perpetual, fixed
	MaxMintsGlobal int64  // total mint times allowed
	MintAmount     int64  // mint amount once
	MaxMints       int64  // # In the fixed mode there is a max number of mints allowed and then no more, only used when mintMode="fixed"
	MaxSupply      int64  // total supply = MaxMintsGlobal*MintAmount
	MintHeight     int64  // start mint height
	MintedTimes    int64  // record minted times
	MintBitworkc   string
	MintBitworkr   string
	Bitworkc       string
	Bitworkr       string
	// Meta           *witness.Meta
	Md           string // emu:"", "0", "1"
	Bv           string // mint_info['$mint_bitwork_vec'] = bv
	Bci          string // mint_info['$mint_bitworkc_inc'] = bci
	Bri          string // mint_info['$mint_bitworkr_inc'] = bri
	Bcs          int64  // mint_info['$mint_bitworkc_start'] = bcs
	Brs          int64  // mint_info['$mint_bitworkr_start'] = brs
	Maxg         int64
	CommitHeight int64
}

func (*GlobalDistributedFt) TableName() string {
	return globalDistributedFtTableName
}

func (*GlobalDistributedFt) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(globalDistributedFtTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&GlobalDistributedFt{})
	assert.Nil(nil, err)
}

func (*GlobalDistributedFt) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(GlobalDistributedFt{})
	assert.Nil(nil, err)
}
