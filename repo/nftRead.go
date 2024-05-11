package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) NftUTXOsByLocationID(locationID string) ([]*postsql.UTXONftInfo, error) {
	// read from cache first, when this txID is in TxCache, this UTXONftInfo must in UserNftInfoCache; otherwise, this UTXONftInfo is not exist
	// entities, ok := m.UserNftInfoCache[locationID]
	// if ok {
	// 	return entities, nil
	// }

	var entity []*postsql.UTXONftInfo
	dbTx := m.Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	var res []*postsql.UTXONftInfo
	for _, UTXO := range entity {
		res = append(res, &postsql.UTXONftInfo{
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
	var entity *postsql.UTXONftInfo
	dbTx := m.Where("atomicals_id = ?", parentRealmAtomicalsID).Find(&entity)
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
	dbTx := m.Where("atomicals_id = ?", parentContainerAtomicalsID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (m *Postgres) NftRealmByNameHasExist(realmName string) (bool, error) {
	if _, exist := m.RealmCache[realmName]; !exist {
		return false, nil
	}
	return true, nil

	// var entities []*postsql.UTXONftInfo
	// dbTx := m.Where("realm_name = ?", realmName).First(&entities)
	// if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
	// 	return false, dbTx.Error
	// }
	// if dbTx.RowsAffected == 0 {
	// 	return false, nil
	// }
	// if len(entities) == 0 {
	// 	return false, nil
	// }
	// return true, nil
}

func (m *Postgres) NftSubRealmByNameHasExist(realmName, subRealm string) (bool, error) {
	if _, exist := m.RealmCache[realmName]; !exist {
		return false, nil
	}
	if _, exist := m.RealmCache[realmName][subRealm]; !exist {
		return false, nil
	}
	return true, nil

	// var entities []*postsql.UTXONftInfo
	// dbTx := m.Where("realm_name = ? and sub_realm_name = ?", realmName, subRealm).First(&entities)
	// if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
	// 	return false, dbTx.Error
	// }
	// if dbTx.RowsAffected == 0 {
	// 	return false, nil
	// }
	// return true, nil
}

func (m *Postgres) NftContainerByNameHasExist(containerName string) (bool, error) {
	if _, exist := m.ContainerCache[containerName]; !exist {
		return false, nil
	}
	return true, nil
	// var entities []*postsql.UTXONftInfo
	// dbTx := m.Where("container_name = ?", containerName).First(&entities)
	// if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
	// 	return false, dbTx.Error
	// }
	// if dbTx.RowsAffected == 0 {
	// 	return false, nil
	// }
	// return true, nil
}

func (m *Postgres) ContainerItemByNameHasExist(containerName, itemID string) (bool, error) {
	if _, exist := m.ContainerCache[containerName]; !exist {
		return false, nil
	}
	if _, exist := m.ContainerCache[containerName][itemID]; !exist {
		return false, nil
	}
	return true, nil

	// var entities []*postsql.UTXONftInfo
	// dbTx := m.Where("container_name = ? and dmitem = ?", containerName, itemID).Find(&entities)
	// if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
	// 	return false, dbTx.Error
	// }
	// if dbTx.RowsAffected == 0 {
	// 	return false, nil
	// }
	// return true, nil
}
