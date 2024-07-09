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

func (m *Postgres) InsertBloomFilter(name string, filter *bloom.BloomFilter) error {
	data, err := filter.MarshalJSON()
	if err != nil {
		return err
	}
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return m.Save(&postsql.BloomFilter{
			Name: name,
			Data: data,
		})
	}) + ";"
	return nil
}

func (m *Postgres) UpdateBloomFilter(name string, filter *bloom.BloomFilter) error {
	data, err := filter.MarshalJSON()
	if err != nil {
		return err
	}
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.BloomFilter{}).Where("name = ?", name).Update("data", data)
	}) + ";"
	return nil
}

func (m *Postgres) addRealm(realm string) {
	m.bloomFilter[postsql.NftFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftRealm, realm)))
	m.bloomFilter[postsql.NftFilter].needUpdate = true
}

func (m *Postgres) testRealm(realm string) bool {
	return m.bloomFilter[postsql.NftFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftRealm, realm)))
}

// func (m *BloomFilterMap) AddSubRealm(realm, subRealm string) {
// 	m.Filter["nft"].Add([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftSubRealm, realm, subRealm)))
// }

// func (m *BloomFilterMap) TestSubRealm(realm, subRealm string) bool {
// 	return m.Filter["nft"].Test([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftSubRealm, realm, subRealm)))
// }

func (m *Postgres) addContainer(container string) {
	m.bloomFilter[postsql.NftFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftContainer, container)))
	m.bloomFilter[postsql.NftFilter].needUpdate = true
}

func (m *Postgres) testContainer(container string) bool {
	return m.bloomFilter[postsql.NftFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftContainer, container)))
}

// func (m *BloomFilterMap) AddItem(container, item string) {
// 	m.Filter["nft"].Add([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftItem, container, item)))
// }

// func (m *BloomFilterMap) TestItem(container, item string) bool {
// 	return m.Filter["nft"].Test([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftItem, container, item)))
// }

func (m *Postgres) addDistributedFt(ftName string) {
	m.bloomFilter[postsql.FtFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDistributedFt, ftName)))
	m.bloomFilter[postsql.FtFilter].needUpdate = true
}

func (m *Postgres) testDistributedFt(ftName string) bool {
	return m.bloomFilter[postsql.FtFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDistributedFt, ftName)))
}

func (m *Postgres) addDirectFt(ftName string) {
	m.bloomFilter[postsql.FtFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, ftName)))
	m.bloomFilter[postsql.FtFilter].needUpdate = true
}

func (m *Postgres) testDirectFt(ftName string) bool {
	return m.bloomFilter[postsql.FtFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, ftName)))
}

func (m *Postgres) addNftLocationID(locationID string) {
	m.bloomFilter[postsql.NftLocationFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
	m.bloomFilter[postsql.NftLocationFilter].needUpdate = true
}

func (m *Postgres) testNftLocationID(locationID string) bool {
	return m.bloomFilter[postsql.NftLocationFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

func (m *Postgres) addFtLocationID(locationID string) {
	m.bloomFilter[postsql.FtLocationFilter].filter.Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
	m.bloomFilter[postsql.FtLocationFilter].needUpdate = true
}

func (m *Postgres) testFtLocationID(locationID string) bool {
	return m.bloomFilter[postsql.FtLocationFilter].filter.Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

func (m *Postgres) InitBloomFilter() {
	limit := 1000
	for offset := 0; offset < 17000; offset += limit {
		nfts, err := m.NftUTXOsByID(offset, limit)
		if err != nil {
			panic(err)
		}
		for _, nft := range nfts {
			if nft.RealmName != "" {
				m.addRealm(nft.RealmName)
			}
			if nft.ContainerName != "" {
				m.addContainer(nft.ContainerName)
			}
			m.addNftLocationID(nft.LocationID)
		}
	}
	for offset := 0; offset < 50000; offset += limit {
		fts, err := m.FtUTXOsByID(offset, limit)
		if err != nil {
			panic(err)
		}
		for _, ft := range fts {
			m.addFtLocationID(ft.LocationID)
		}
	}
	for offset := 0; offset < 50000; offset += limit {
		fts, err := m.DistributedFtByID(offset, limit)
		if err != nil {
			panic(err)
		}
		for _, ft := range fts {
			m.addDistributedFt(ft.TickerName)
		}
	}
	for name, v := range m.bloomFilter {
		m.UpdateBloomFilter(name, v.filter)
	}
}
