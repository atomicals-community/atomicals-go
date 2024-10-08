package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) NftUTXOsByUserPK(UserPK string) ([]*postsql.UTXONftInfo, error) {
	var entity []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("user_pk = ?", UserPK).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}

func (m *Postgres) NftUTXOByAtomicalsID(atomicalsID string) (*postsql.UTXONftInfo, error) {
	var entity *postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("atomicals_id = ?", atomicalsID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}

func (m *Postgres) NftUTXOsByLocationID(locationID string) ([]*postsql.UTXONftInfo, error) {
	if !m.testNftLocationID(locationID) {
		return nil, nil
	}
	var entity []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}

func (m *Postgres) NftRealmByName(realmName string) ([]*postsql.UTXONftInfo, error) {
	if !m.testRealm(realmName) {
		return nil, nil
	}
	var entities []*postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("realm_name = ?", realmName).First(&entities)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}

	if len(entities) == 0 {
		return nil, nil
	}
	return entities, nil
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
	if !m.testContainer(containerName) {
		return false, nil
	}
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

func (m *Postgres) LatestItemByContainerName(containerName string) (*postsql.UTXONftInfo, error) {
	var entity *postsql.UTXONftInfo
	dbTx := m.Model(postsql.UTXONftInfo{}).Where("container_name = ?", containerName).Order("id DESC").Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return nil, nil
}
