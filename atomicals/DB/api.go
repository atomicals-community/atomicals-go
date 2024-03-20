package db

import (
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/btcsync"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	// atomicals-core current location
	CurrentHeitht() (int64, error)
	CurrentLocation() (int64, string, error)
	UpdateLocation(int64, string) error

	// NftUTXOsByAtomicalsID and NftUTXOsByLocationID are two index map from one table
	// NftUTXOsByAtomicalsID: indexed by AtomicalsID, NftUTXOsByLocationID: indexed by LocationID
	NftUTXOsByAtomicalsID(atomicalsID string) ([]*UserNftInfo, error)
	NftUTXOsByLocationID(locationID string) ([]*UserNftInfo, error)
	InsertNftUTXO(UTXO *UserNftInfo) error
	TransferNftUTXO(oldLocationID, newLocationID, newUserPk string) error
	InsertGlobalNft(nftType int64, name, subName string) error

	// an golbal map to check if realm nft has been registed
	ParentRealmHasExist(parentRealmAtomicalsID string) (string, error)
	NftRealmByNameHasExist(realmName string) (bool, error)
	NftSubRealmByNameHasExist(realmName, subRealm string) (bool, error)

	// an golbal map to check if container nft has been registed
	ParentContainerHasExist(parentContainerAtomicalsID string) (string, error)
	NftContainerByNameHasExist(containerName string) (bool, error)
	ContainerItemByNameHasExist(container, item string) (bool, error)

	// FtUTXOsByLocationID: fts indexed by LocationID (operation: dmt, mint ft)
	FtUTXOsByLocationID(locationID string) ([]*UserFtInfo, error)
	InsertFtUTXO(UTXO *UserFtInfo) error
	DeleteFtUTXO(locationID string) error

	// DistributedFt: fts indexed by tickerName (operation: dft, deploy ft)
	DistributedFtByName(tickerName string) (*DistributedFtInfo, error)
	InsertDistributedFt(ft *DistributedFtInfo) error
	UpdateDistributedFtAmount(tickerName string, mintTimes int64) error

	// DirectFt: Direct fts indexed by LocationID (operation: dft, mint direct ft)
	DirectFtByName(tickerName string) (bool, error)
	InsertDirectFt(tickerName string) error
}

func NewSqlDB(sqlDNS string, btcClient *btcsync.BtcSync) DB {
	DB, err := gorm.Open(postgres.Open(sqlDNS))
	if err != nil {
		panic(err)
	}
	m := &Postgres{
		DB:               DB,
		SQLRaw:           "",
		UserNftInfoCache: make(map[string][]*UserNftInfo, 0),
		UserFtInfoCache:  make(map[string][]*UserFtInfo, 0),
	}
	_, err = m.CurrentHeitht()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := m.UpdateLocation(common.ATOMICALS_ACTIVATION_HEIGHT-1, ""); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return m
}

// NewMemoryDB will be supported later
// func NewMemoryDB() DB {
// 	return &Memory{
// 		height:                common.ATOMICALS_ACTIVATION_HEIGHT - 1,
// 		txID:                  "",
// 		nftUTXOsByAtomicalsID: make(map[string][]*UserNftInfo, 0),
// 		nftUTXOsByLocationID:  make(map[string][]*UserNftInfo, 0),
// 		globalNftRealmMap:     make(map[string]map[string]bool, 0),
// 		globalNftContainerMap: make(map[string]map[string]bool, 0),

// 		ftUTXOs:                make(map[string][]*UserFtInfo, 0),
// 		globalDistributedFtMap: make(map[string]*DistributedFtInfo, 0),
// 		globalDirectFtMap:      make(map[string]bool, 0),
// 	}
// }
