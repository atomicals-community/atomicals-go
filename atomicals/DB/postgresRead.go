package db

import (
	"github.com/atomicals-core/atomicals/DB/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) CurrentHeitht() (int64, error) {
	entity := &postsql.Location{}
	dbTx := m.Find(&entity)
	if dbTx.Error != nil {
		return -1, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return -1, gorm.ErrRecordNotFound
	}
	return entity.Height, nil
}

func (m *Postgres) CurrentLocation() (int64, string, error) {
	entity := &postsql.Location{}
	dbTx := m.Take(&entity)
	if dbTx.Error != nil {
		return -1, "", dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return -1, "", gorm.ErrRecordNotFound
	}
	return entity.Height, entity.TxID, nil
}

func (m *Postgres) NftUTXOsByAtomicalsID(atomicalsID string) ([]*UserNftInfo, error) {
	var entity []*postsql.UserNftInfo
	dbTx := m.Where("atomicals_id = ?", atomicalsID).Find(&entity)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	var res []*UserNftInfo
	for _, UTXO := range entity {
		res = append(res, &UserNftInfo{
			UserPk:                     UTXO.UserPk,
			AtomicalsID:                UTXO.AtomicalsID,
			LocationID:                 UTXO.LocationID,
			RealmName:                  UTXO.RealmName,
			SubRealmName:               UTXO.SubRealmName,
			ClaimType:                  UTXO.ClaimType,
			ParentRealmAtomicalsID:     UTXO.ParentRealmAtomicalsID,
			ContainerName:              UTXO.ContainerName,
			Dmitem:                     UTXO.Dmitem,
			ParentContainerAtomicalsID: UTXO.ParentContainerAtomicalsID,
			Nonce:                      UTXO.Nonce,
			Time:                       UTXO.Time,
		})
	}
	return res, nil
}

func (m *Postgres) NftUTXOsByLocationID(locationID string) ([]*UserNftInfo, error) {
	// read from cache first, when this txID is in TxCache, this UserNftInfo must in UserNftInfoCache; otherwise, this UserNftInfo is not exist
	entities, ok := m.UserNftInfoCache[locationID]
	if ok {
		return entities, nil
	}

	var entity []*postsql.UserNftInfo
	dbTx := m.Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	var res []*UserNftInfo
	for _, UTXO := range entity {
		res = append(res, &UserNftInfo{
			UserPk:                     UTXO.UserPk,
			AtomicalsID:                UTXO.AtomicalsID,
			LocationID:                 UTXO.LocationID,
			RealmName:                  UTXO.RealmName,
			SubRealmName:               UTXO.SubRealmName,
			ClaimType:                  UTXO.ClaimType,
			ParentRealmAtomicalsID:     UTXO.ParentRealmAtomicalsID,
			ContainerName:              UTXO.ContainerName,
			Dmitem:                     UTXO.Dmitem,
			ParentContainerAtomicalsID: UTXO.ParentContainerAtomicalsID,
			Nonce:                      UTXO.Nonce,
			Time:                       UTXO.Time,
		})
	}
	return res, nil
}

func (m *Postgres) ParentRealmHasExist(parentRealmAtomicalsID string) (string, error) {
	var entity *postsql.UserNftInfo
	dbTx := m.Where("atomicals_id = ?", parentRealmAtomicalsID).Find(&entity)
	if dbTx.Error != nil {
		return "", dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return "", nil
	}
	return entity.RealmName, nil
}

func (m *Postgres) ParentContainerHasExist(parentContainerAtomicalsID string) (string, error) {
	var entity *postsql.UserNftInfo
	dbTx := m.Where("atomicals_id = ?", parentContainerAtomicalsID).Find(&entity)
	if dbTx.Error != nil {
		return "", dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return "", nil
	}
	return entity.ContainerName, nil
}

func (m *Postgres) NftRealmByNameHasExist(realmName string) (bool, error) {
	var entities []*postsql.GlobalNft
	dbTx := m.Where("nft_type = ? and name = ?", TypeNftRealm, realmName).Find(&entities)
	if dbTx.Error != nil {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	if len(entities) == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) NftSubRealmByNameHasExist(realmName, subRealm string) (bool, error) {
	var entities []*postsql.GlobalNft
	dbTx := m.Where("nft_type = ? and name = ? and sub_name = ?", TypeNftSubRealm, realmName, subRealm).Find(&entities)
	if dbTx.Error != nil {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) NftContainerByNameHasExist(containerName string) (bool, error) {
	var entities []*postsql.GlobalNft
	dbTx := m.Where("nft_type = ? and name = ?", TypeNftContainer, containerName).Find(&entities)
	if dbTx.Error != nil {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) ContainerItemByNameHasExist(containerName, itemID string) (bool, error) {
	var entities []*postsql.GlobalNft
	dbTx := m.Where("nft_type = ? and name = ? and sub_name = ?", TypeNftItem, containerName, itemID).Find(&entities)
	if dbTx.Error != nil {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) FtUTXOsByLocationID(locationID string) ([]*UserFtInfo, error) {
	// read from cache first, when this txID is in TxCache, this UserNftInfo must in UserNftInfoCache; otherwise, this UserNftInfo is not exist
	entities, ok := m.UserFtInfoCache[locationID]
	if ok {
		return entities, nil
	}

	var entity []*postsql.UserFtInfo
	dbTx := m.Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	var res []*UserFtInfo
	for _, UTXO := range entity {
		res = append(res, &UserFtInfo{
			UserPk:        UTXO.UserPk,
			AtomicalsID:   UTXO.AtomicalsID,
			LocationID:    UTXO.LocationID,
			MintTicker:    UTXO.MintTicker,
			Nonce:         UTXO.Nonce,
			Time:          UTXO.Time,
			Amount:        UTXO.Amount,
			Type:          UTXO.Type,
			Subtype:       UTXO.Subtype,
			RequestTicker: UTXO.RequestTicker,
			MaxSupply:     UTXO.MaxSupply,
			MintAmount:    UTXO.MintAmount,
			MintHeight:    UTXO.MintHeight,
			MaxMints:      UTXO.MaxMints,
		})
	}
	return res, nil
}

func (m *Postgres) DistributedFtByName(tickerName string) (*DistributedFtInfo, error) {
	var entity *postsql.GlobalDistributedFt
	dbTx := m.Where("ticker_name = ?", tickerName).Find(&entity)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return &DistributedFtInfo{
		AtomicalsID:    entity.AtomicalsID,
		TickerName:     entity.TickerName,
		Type:           entity.Type,
		Subtype:        entity.Subtype,
		MintMode:       entity.MintMode,
		MaxMintsGlobal: entity.MaxMintsGlobal,
		MintAmount:     entity.MintAmount,
		MaxMints:       entity.MaxMints,
		MaxSupply:      entity.MaxSupply,
		MintHeight:     entity.MintHeight,
		MintedTimes:    entity.MintedTimes,
		Md:             entity.Md,
		Bv:             entity.Bv,
		Bci:            entity.Bci,
		Bri:            entity.Bri,
		Bcs:            entity.Bcs,
		Brs:            entity.Brs,
		Maxg:           entity.Maxg,
		CommitHeight:   entity.CommitHeight,
	}, nil
}

func (m *Postgres) InsertDistributedFt(entity *DistributedFtInfo) error {
	dbTx := m.Save(&postsql.GlobalDistributedFt{
		AtomicalsID:    entity.AtomicalsID,
		TickerName:     entity.TickerName,
		Type:           entity.Type,
		Subtype:        entity.Subtype,
		MintMode:       entity.MintMode,
		MaxMintsGlobal: entity.MaxMintsGlobal,
		MintAmount:     entity.MintAmount,
		MaxMints:       entity.MaxMints,
		MaxSupply:      entity.MaxSupply,
		MintHeight:     entity.MintHeight,
		MintedTimes:    entity.MintedTimes,
		Md:             entity.Md,
		Bv:             entity.Bv,
		Bci:            entity.Bci,
		Bri:            entity.Bri,
		Bcs:            entity.Bcs,
		Brs:            entity.Brs,
		Maxg:           entity.Maxg,
		CommitHeight:   entity.CommitHeight,
	})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) DirectFtByName(tickerName string) (bool, error) {
	var entity *postsql.GlobalDirectFt
	dbTx := m.Where("ticker_name = ?", tickerName).Find(&entity)
	if dbTx.Error != nil {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
