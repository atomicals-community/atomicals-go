package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) NftUTXOsByLocationID(locationID string) ([]*postsql.UTXONftInfo, error) {
	var entity []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (m *Postgres) ParentRealmHasExist(parentRealmAtomicalsID string) (string, error) {
	var entity *postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("atomicals_id = ?", parentRealmAtomicalsID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return "", dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return "", nil
	}
	return entity.RealmName, nil
}

func (m *Postgres) ParentContainerHasExist(parentContainerAtomicalsID string) (*postsql.UTXONftInfo, error) {
	var entity *postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("atomicals_id = ?", parentContainerAtomicalsID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (m *Postgres) NftRealmByNameHasExist(realmName string) (bool, error) {
	var entities []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("realm_name = ?", realmName).First(&entities)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
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

func (m *Postgres) NftSubRealmByNameHasExist(parentRealmAtomicalsID, subRealm string) (bool, error) {
	var entities []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("parent_realm_atomicals_id = ? and sub_realm_name = ?", parentRealmAtomicalsID, subRealm).First(&entities)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) NftContainerByNameHasExist(containerName string) (bool, error) {
	var entities []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("container_name = ?", containerName).First(&entities)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (m *Postgres) ContainerItemByNameHasExist(containerName, itemID string) (bool, error) {
	var entities []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("container_name = ? and dmitem = ?", containerName, itemID).Find(&entities)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return false, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
