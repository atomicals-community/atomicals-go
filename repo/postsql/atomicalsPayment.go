package postsql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const PaymentTableName = "atomicals_payment"

type PaymentInfo struct {
	gorm.Model
	Height      int64
	AtomicalsID string `gorm:"index"` // txID _ VOUT_EXPECT_OUTPUT_INDEX when be minted
	LocationID  string `gorm:"index"` // txID_voutIndex updated after being transfered
	Payment     string
	PaymentStr  string
}

func (*PaymentInfo) TableName() string {
	return PaymentTableName
}

func (*PaymentInfo) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(PaymentTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&PaymentInfo{})
	assert.Nil(nil, err)
}

func (*PaymentInfo) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(PaymentInfo{})
	assert.Nil(nil, err)
}
