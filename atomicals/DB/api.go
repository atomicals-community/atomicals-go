package db

import (
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
	InsertNftUTXOByAtomicalsID(UTXO *UserNftInfo) error
	NftUTXOsByLocationID(locationID string) ([]*UserNftInfo, error)
	InsertNftUTXOByLocationID(UTXO *UserNftInfo) error
	DeleteNftUTXOByLocationID(locationID string) error

	// an golbal map to check if realm nft has been registed
	ParentRealmHasExist(parentRealmAtomicalsID string) (string, error)
	NftRealmByName(realmName string) (map[string]bool, error)
	NftSubRealmByName(realmName, subRealm string) (bool, error)
	InsertRealm(realmName string) error
	InsertSubRealm(realmName, subRealm string) error

	// an golbal map to check if container nft has been registed
	ParentContainerHasExist(parentContainerAtomicalsID string) (string, error)
	NftContainerByName(containerName string) (map[string]bool, error)
	InsertContainer(containerName string) error
	InsertItemInContainer(containerName, itemID string) error

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

func NewMemoryDB(height int64, txID string) DB {
	return &Memory{
		height:                height,
		txID:                  txID,
		nftUTXOsByAtomicalsID: make(map[string][]*UserNftInfo, 0),
		nftUTXOsByLocationID:  make(map[string][]*UserNftInfo, 0),
		globalNftRealmMap:     make(map[string]map[string]bool, 0),
		globalNftContainerMap: make(map[string]map[string]bool, 0),

		ftUTXOs:                make(map[string][]*UserFtInfo, 0),
		globalDistributedFtMap: make(map[string]*DistributedFtInfo, 0),
		globalDirectFtMap:      make(map[string]bool, 0),
	}
}

func NewSqlDB(DB *gorm.DB, height int64, txID string) DB {
	m := &Postgres{
		DB,
	}
	if err := m.UpdateLocation(height, txID); err != nil {
		panic(err)
	}
	return m
}
