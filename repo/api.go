package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	// atomicals-core current location
	CurrentLocation() (*postsql.Location, error)
	UpdateCurrentHeightAndExecAllSql(int64, int64) error

	// nft read
	NftUTXOsByLocationID(locationID string) ([]*postsql.UTXONftInfo, error)
	ParentRealmHasExist(parentRealmAtomicalsID string) (string, error)
	NftRealmByNameHasExist(realmName string) (bool, error)
	NftSubRealmByNameHasExist(realmName, subRealm string) (bool, error)
	ParentContainerHasExist(parentContainerAtomicalsID string) (*postsql.UTXONftInfo, error)
	NftContainerByNameHasExist(containerName string) (bool, error)
	ContainerItemByNameHasExist(container, item string) (bool, error)
	// nft write
	InsertNftUTXO(UTXO *postsql.UTXONftInfo) error
	TransferNftUTXO(oldLocationID, newLocationID, newUserPk string) error

	// ft read
	FtUTXOsByLocationID(locationID string) ([]*postsql.UTXOFtInfo, error)
	DistributedFtByName(tickerName string) (*postsql.GlobalDistributedFt, error)
	// ft write
	InsertFtUTXO(UTXO *postsql.UTXOFtInfo) error
	DeleteFtUTXO(locationID string) error
	InsertDistributedFt(ft *postsql.GlobalDistributedFt) error
	UpdateDistributedFtAmount(tickerName string, mintTimes int64) error

	// mod
	InsertMod(mod *postsql.ModInfo) error
	Mod(atomicalsID string) (*postsql.ModInfo, error)

	// btc
	InsertBtcTx(btcTx *postsql.BtcTx) error
	BtcTx(txID string) (*postsql.BtcTx, error)
	BtcTxHeight(txID string) (int64, error)
}

func NewSqlDB(sqlDNS string) DB {
	DB, err := gorm.Open(postgres.Open(sqlDNS))
	if err != nil {
		panic(err)
	}
	return &Postgres{
		DB:     DB,
		SQLRaw: "",
	}
}
