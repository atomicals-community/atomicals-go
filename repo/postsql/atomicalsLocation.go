package postsql

import (
	"github.com/atomicals-go/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const locationTableName = "atomicals_location"

type Location struct {
	gorm.Model
	Owner   string `gorm:"uniqueindex" json:"owner"`
	Height  int64
	TxIndex int64
}

func (*Location) TableName() string {
	return locationTableName
}

func (*Location) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(locationTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&Location{})
	assert.Nil(nil, err)
	dbTx := db.Save(&Location{Owner: "atomicals", Height: utils.ATOMICALS_ACTIVATION_HEIGHT, TxIndex: -1})
	assert.Nil(nil, dbTx.Error)
}

func (*Location) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(Location{})
	assert.Nil(nil, err)
}
