package postsql

import (
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	NftFilter         = "nft"
	FtFilter          = "ft"
	NftLocationFilter = "nft_locationID"
	FtLocationFilter  = "ft_locationID"
)

const BloomFilterTableName = "bloomFilter"

type BloomFilter struct {
	gorm.Model
	Name string `gorm:"uniqueindex"`
	Data []byte
}

func (*BloomFilter) TableName() string {
	return BloomFilterTableName
}

func (*BloomFilter) Init(db *gorm.DB) {
	var err error
	dmodel := newDefaultModel(BloomFilterTableName, db)
	err = dmodel.DropTable()
	assert.Nil(nil, err)
	err = dmodel.CreateTable(&BloomFilter{})
	assert.Nil(nil, err)

	filter := bloom.NewWithEstimates(10000, 0.01)
	data, err := filter.MarshalJSON()
	assert.Nil(nil, err)
	dbTx := db.Save(&BloomFilter{Name: NftFilter, Data: data})
	assert.Nil(nil, dbTx.Error)
	filter = bloom.NewWithEstimates(10000, 0.01)
	data, err = filter.MarshalJSON()
	assert.Nil(nil, err)
	dbTx = db.Save(&BloomFilter{Name: FtFilter, Data: data})
	assert.Nil(nil, dbTx.Error)

	filter = bloom.NewWithEstimates(40000, 0.01)
	data, err = filter.MarshalJSON()
	assert.Nil(nil, err)
	dbTx = db.Save(&BloomFilter{Name: NftLocationFilter, Data: data})
	assert.Nil(nil, dbTx.Error)

	filter = bloom.NewWithEstimates(80000, 0.01)
	data, err = filter.MarshalJSON()
	assert.Nil(nil, err)
	dbTx = db.Save(&BloomFilter{Name: FtLocationFilter, Data: data})
	assert.Nil(nil, dbTx.Error)
}

func (*BloomFilter) AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(BloomFilter{})
	assert.Nil(nil, err)
}
