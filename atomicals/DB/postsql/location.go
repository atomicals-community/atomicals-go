package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const locationTableName = "atomicals_location"

type Location struct {
	gorm.Model
	Owner  string `gorm:"uniqueindex" json:"owner"`
	Height int64
	TxID   string
}

func (*Location) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(locationTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&Location{})
	assert.Nil(nil, err)
}

func (*Location) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(Location{})
	assert.Nil(nil, err)
}
