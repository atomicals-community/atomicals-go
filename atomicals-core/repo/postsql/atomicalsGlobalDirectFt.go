package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const globalDirectFtTableName = "atomicals_global_direct_tf"

type GlobalDirectFt struct {
	gorm.Model
	TickerName string `gorm:"uniqueindex"`
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
