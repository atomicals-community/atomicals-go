package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source api.go -destination api_mock.go -package repo
type DB interface {
	// location
	Location() (*postsql.Location, error)
	ExecAllSql(blockHeight, txIndex int64, txID, operation string) error

	// nft read
	NftUTXOsByUserPK(UserPK string) ([]*postsql.UTXONftInfo, error)
	NftUTXOsByLocationID(locationID string) ([]*postsql.UTXONftInfo, error)
	ParentRealmHasExist(parentRealmAtomicalsID string) (string, error)
	NftRealmByNameHasExist(realmName string) (bool, error)
	NftSubRealmByNameHasExist(parentRealmAtomicalsID, subRealm string) (bool, error)
	ParentContainerHasExist(parentContainerAtomicalsID string) (*postsql.UTXONftInfo, error)
	NftContainerByNameHasExist(containerName string) (bool, error)
	ContainerItemByNameHasExist(container, item string) (bool, error)
	// nft write
	InsertNftUTXO(UTXO *postsql.UTXONftInfo) error
	UpdateNftUTXO(UTXO *postsql.UTXONftInfo) error

	// ft read
	FtUTXOsByUserPK(UserPK string) ([]*postsql.UTXOFtInfo, error)
	FtUTXOsByLocationID(locationID string) ([]*postsql.UTXOFtInfo, error)
	DistributedFtByName(tickerName string) (*postsql.GlobalDistributedFt, error)
	// ft write
	InsertFtUTXO(UTXO *postsql.UTXOFtInfo) error
	DeleteFtUTXO(locationID string) error
	InsertDistributedFt(ft *postsql.GlobalDistributedFt) error
	UpdateDistributedFt(ft *postsql.GlobalDistributedFt) error
	InsertDirectFtUTXO(entity *postsql.GlobalDirectFt) error

	// mod
	InsertMod(mod *postsql.ModInfo) error
	Mod(atomicalsID string) (*postsql.ModInfo, error)

	// btc
	InsertBtcTx(btcTx *postsql.BtcTx) error
	BtcTx(txID string) (*postsql.BtcTx, error)
	BtcTxHeight(txID string) (int64, error)
	DeleteBtcTxUntil(blockHeight int64) error

	// BloomFilter
	InsertBloomFilter(name string, filter *bloom.BloomFilter) error
	UpdateBloomFilter(name string, filter *bloom.BloomFilter) error
	BloomFilter() (map[string]*bloomFilterInfo, error)

	// sync nft/ft asset with locationID
	AddLocationIDIntoChannel(locationID string)
	GetNftUTXOFromChannel(locationID string) ([]*postsql.UTXONftInfo, error)
	GetFtUTXOFromChannel(locationID string) ([]*postsql.UTXOFtInfo, error)
}

func NewSqlDB(sqlDNS string) DB {
	DB, err := gorm.Open(postgres.Open(sqlDNS))
	if err != nil {
		panic(err)
	}
	s := &Postgres{
		DB:                DB,
		locationIDChannel: make(chan string, 10),
	}
	s.bloomFilter, err = s.BloomFilter()
	if err != nil {
		panic(err)
	}
	go s.fetchUTXO()
	return s
}
