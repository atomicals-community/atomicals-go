package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:generate mockgen -source api.go -destination api_mock.go -package repo
type DB interface {
	// location
	Location() (*postsql.Location, error)

	// btc
	AtomicalsTx(txID string) (*postsql.AtomicalsTx, error)
	AtomicalsTxByHeight(height int64) ([]*postsql.AtomicalsTx, error)
	AtomicalsTxHeight(txID string) (int64, error)

	// nft read
	NftUTXOsByUserPK(UserPK string) ([]*postsql.UTXONftInfo, error)
	NftUTXOByAtomicalsID(atomicalsID string) (*postsql.UTXONftInfo, error)
	NftUTXOsByLocationID(locationID string) ([]*postsql.UTXONftInfo, error)
	NftRealmByName(realmName string) ([]*postsql.UTXONftInfo, error)
	NftSubRealmByNameHasExist(parentRealmAtomicalsID, subRealm string) (bool, error)
	NftContainerByNameHasExist(containerName string) (bool, error)
	ContainerItemByNameHasExist(container, item string) (bool, error)
	LatestItemByContainerName(container string) (*postsql.UTXONftInfo, error)

	// ft read
	FtUTXOsByUserPK(UserPK string) ([]*postsql.UTXOFtInfo, error)
	FtUTXOsByLocationID(locationID string) ([]*postsql.UTXOFtInfo, error)
	DistributedFtByName(tickerName string) (*postsql.GlobalDistributedFt, error)
	DirectFtByName(tickerName string) (*postsql.GlobalDirectFt, error)

	// mod
	ModHistory(atomicalsID string, height int64) ([]*postsql.ModInfo, error)

	UpdateDB(currentHeight, currentTxIndex int64, txID string, data *AtomicaslData) error

	PostgresDB() *Postgres
}

func NewSqlDB(sqlDNS string) DB {
	DB, err := gorm.Open(postgres.Open(sqlDNS))
	if err != nil {
		panic(err)
	}
	s := &Postgres{
		DB: DB,
	}
	s.bloomFilter, err = s.BloomFilter()
	if err != nil {
		panic(err)
	}
	return s
}
