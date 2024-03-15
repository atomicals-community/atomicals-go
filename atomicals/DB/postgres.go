package db

import (
	"github.com/atomicals-core/atomicals/DB/postsql"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
}

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

func (m *Postgres) UpdateLocation(height int64, txID string) error {
	entity := &postsql.Location{}
	dbTx := m.Take(&entity)
	if dbTx.Error != nil {
		return dbTx.Error
	}
	entity.Height = height
	entity.TxID = txID
	dbTX := m.Save(entity)
	if dbTX.Error != nil {
		return dbTx.Error
	}
	return nil
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
			NftType:                    UTXO.NftType,
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

func (m *Postgres) InsertNftUTXOByAtomicalsID(UTXO *UserNftInfo) error {
	dbTx := m.Save(&postsql.UserNftInfo{
		UserPk:                     UTXO.UserPk,
		AtomicalsID:                UTXO.AtomicalsID,
		LocationID:                 UTXO.LocationID,
		NftType:                    UTXO.NftType,
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
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) NftUTXOsByLocationID(locationID string) ([]*UserNftInfo, error) {
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
			NftType:                    UTXO.NftType,
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

func (m *Postgres) InsertNftUTXOByLocationID(UTXO *UserNftInfo) error {
	dbTx := m.Save(&postsql.UserNftInfo{
		UserPk:                     UTXO.UserPk,
		AtomicalsID:                UTXO.AtomicalsID,
		LocationID:                 UTXO.LocationID,
		NftType:                    UTXO.NftType,
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
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) DeleteNftUTXOByLocationID(locationID string) error {
	dbTx := m.Where("loaction_id = ?", locationID).Delete(&postsql.UserNftInfo{})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil
	}
	return nil
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

func (m *Postgres) NftRealmByName(realmName string) (map[string]bool, error) {
	var entities []*postsql.GlobalNftRealm
	dbTx := m.Where("realm_name = ?", realmName).Find(&entities)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	res := make(map[string]bool, 0)
	for _, v := range entities {
		res[v.SubRealmName] = true
	}
	return res, nil
}

func (m *Postgres) NftSubRealmByName(realmName, subRealm string) (bool, error) {
	var entity *postsql.GlobalNftRealm
	dbTx := m.Where("realm_name = ? and subRealm = ?", realmName, subRealm).Find(&entity)
	if dbTx.Error != nil {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) InsertRealm(realmName string) error {
	dbTx := m.Save(&postsql.GlobalNftRealm{
		RealmName:    realmName,
		SubRealmName: "",
	})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) InsertSubRealm(realmName, subRealm string) error {
	dbTx := m.Save(&postsql.GlobalNftRealm{
		RealmName:    realmName,
		SubRealmName: subRealm,
	})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) NftContainerByName(containerName string) (map[string]bool, error) {
	var entities []*postsql.GlobalNftContainer
	dbTx := m.Where("container_name = ?", containerName).Find(&entities)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	res := make(map[string]bool, 0)
	for _, v := range entities {
		res[v.ContainerName] = true
	}
	return res, nil
}

func (m *Postgres) InsertContainer(containerName string) error {
	dbTx := m.Save(&postsql.GlobalNftContainer{
		ContainerName: containerName,
		ItemID:        "",
	})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) InsertItemInContainer(containerName, itemID string) error {
	dbTx := m.Save(&postsql.GlobalNftRealm{
		RealmName:    containerName,
		SubRealmName: itemID,
	})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) FtUTXOsByLocationID(locationID string) ([]*UserFtInfo, error) {
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

func (m *Postgres) InsertFtUTXO(UTXO *UserFtInfo) error {
	dbTx := m.Save(&postsql.UserFtInfo{
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
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *Postgres) DeleteFtUTXO(locationID string) error {
	dbTx := m.Where("loaction_id = ?", locationID).Delete(&postsql.UserFtInfo{})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil
	}
	return nil
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
		Ticker:         entity.Ticker,
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
		Ticker:         entity.Ticker,
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

func (m *Postgres) UpdateDistributedFtAmount(tickerName string, mintTimes int64) error {
	dbTx := m.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", tickerName).Update("minted_times", mintTimes)
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

func (m *Postgres) InsertDirectFt(tickerName string) error {
	dbTx := m.Save(&postsql.GlobalDirectFt{
		TickerName: tickerName,
	})
	if dbTx.Error != nil {
		return dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
