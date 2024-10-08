package repo

import (
	"fmt"

	"github.com/atomicals-go/repo/postsql"
	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/gorm"
)

type bloomFilterInfo struct {
	filter     *bloom.BloomFilter
	needUpdate bool
}

// ft
func (m *Postgres) addFtLocationID(locationID string) {
	m.bloomFilter[postsql.FtLocationFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
	m.bloomFilter[postsql.FtLocationFilter].needUpdate = true
}

func (m *Postgres) testFtLocationID(locationID string) bool {
	return m.bloomFilter[postsql.FtLocationFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

// func (m *Postgres) addDistributedFt(ftName string) {
// 	m.bloomFilter[postsql.FtFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDistributedFt, ftName)))
// 	m.bloomFilter[postsql.FtFilter].needUpdate = true
// }

// func (m *Postgres) testDistributedFt(ftName string) bool {
// 	return m.bloomFilter[postsql.FtFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDistributedFt, ftName)))
// }

// nft
func (m *Postgres) addRealm(realm string) {
	m.bloomFilter[postsql.NftFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftRealm, realm)))
	m.bloomFilter[postsql.NftFilter].needUpdate = true
}

func (m *Postgres) testRealm(realm string) bool {
	return m.bloomFilter[postsql.NftFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftRealm, realm)))
}

func (m *Postgres) addContainer(container string) {
	m.bloomFilter[postsql.NftFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftContainer, container)))
	m.bloomFilter[postsql.NftFilter].needUpdate = true
}

func (m *Postgres) testContainer(container string) bool {
	return m.bloomFilter[postsql.NftFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftContainer, container)))
}

func (m *Postgres) addNftLocationID(locationID string) {
	m.bloomFilter[postsql.NftLocationFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
	m.bloomFilter[postsql.NftLocationFilter].needUpdate = true
}

func (m *Postgres) testNftLocationID(locationID string) bool {
	return m.bloomFilter[postsql.NftLocationFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

func (m *Postgres) BloomFilter() (map[string]*bloomFilterInfo, error) {
	entities := make([]*postsql.BloomFilter, 0)
	dbTx := m.Order("id desc").Find(&entities)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	filterMap := make(map[string]*bloomFilterInfo, 0)
	for _, v := range entities {
		filter := bloom.NewWithEstimates(10000, 0.01)
		filter.UnmarshalJSON(v.Data)
		filterMap[v.Name] = &bloomFilterInfo{
			filter:     filter,
			needUpdate: false,
		}
	}
	return filterMap, nil
}
