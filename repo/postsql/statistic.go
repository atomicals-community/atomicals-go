package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const statisticTxTableName = "statistic_tx"

type StatisticTx struct {
	gorm.Model
	BlockHeight int64 `gorm:"uniqueindex"`
	Data        string
	Description string
}

func (*StatisticTx) TableName() string {
	return statisticTxTableName
}

func (*StatisticTx) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(statisticTxTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&StatisticTx{})
	assert.Nil(nil, err)
}

func (*StatisticTx) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(StatisticTx{})
	assert.Nil(nil, err)
}
