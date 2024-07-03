package bloomfilter

import (
	"fmt"

	"github.com/atomicals-go/repo/postsql"
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilterMap struct {
	Filter map[string]*bloom.BloomFilter
}

func (m *BloomFilterMap) AddRealm(realm string) {
	m.Filter["nft"].Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftRealm, realm)))
}

func (m *BloomFilterMap) TestRealm(realm string) bool {
	return m.Filter["nft"].Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftRealm, realm)))
}

// func (m *BloomFilterMap) AddSubRealm(realm, subRealm string) {
// 	m.Filter["nft"].Add([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftSubRealm, realm, subRealm)))
// }

// func (m *BloomFilterMap) TestSubRealm(realm, subRealm string) bool {
// 	return m.Filter["nft"].Test([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftSubRealm, realm, subRealm)))
// }

func (m *BloomFilterMap) AddContainer(container string) {
	m.Filter["nft"].Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftContainer, container)))
}

func (m *BloomFilterMap) TestContainer(container string) bool {
	return m.Filter["nft"].Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeNftContainer, container)))
}

// func (m *BloomFilterMap) AddItem(container, item string) {
// 	m.Filter["nft"].Add([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftItem, container, item)))
// }

// func (m *BloomFilterMap) TestItem(container, item string) bool {
// 	return m.Filter["nft"].Test([]byte(fmt.Sprintf("%v_%v_%v", postsql.TypeNftItem, container, item)))
// }

func (m *BloomFilterMap) AddDistributedFt(ftName string) {
	m.Filter["ft"].Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDistributedFt, ftName)))
}

func (m *BloomFilterMap) TestDistributedFt(ftName string) bool {
	return m.Filter["ft"].Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDistributedFt, ftName)))
}

func (m *BloomFilterMap) AddDirectFt(ftName string) {
	m.Filter["ft"].Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, ftName)))
}

func (m *BloomFilterMap) TestDirectFt(ftName string) bool {
	return m.Filter["ft"].Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, ftName)))
}

func (m *BloomFilterMap) AddNftLocationID(locationID string) {
	m.Filter["nft_locationID"].Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

func (m *BloomFilterMap) TestNftLocationID(locationID string) bool {
	return m.Filter["nft_locationID"].Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

func (m *BloomFilterMap) AddFtLocationID(locationID string) {
	m.Filter["ft_locationID"].Add([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}

func (m *BloomFilterMap) TestFtLocationID(locationID string) bool {
	return m.Filter["ft_locationID"].Test([]byte(fmt.Sprintf("%v_%v", postsql.TypeDirectFt, locationID)))
}
