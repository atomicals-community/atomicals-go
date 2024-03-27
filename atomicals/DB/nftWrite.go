package db

import (
	"github.com/atomicals-core/atomicals/DB/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertNftUTXO(UTXO *postsql.UTXONftInfo) error {
	if UTXO.RealmName != "" {
		m.RealmCache[UTXO.RealmName] = true
	}
	if UTXO.ContainerName != "" {
		m.ContainerCache[UTXO.ContainerName] = true
	}

	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.UTXONftInfo{
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
	}) + ";"
	return nil
}

func (m *Postgres) TransferNftUTXO(oldLocationID, newLocationID, newUserPk string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.UTXONftInfo{}).Where("loaction_id = ?", oldLocationID).Updates(map[string]interface{}{"loaction_id": newLocationID, "user_pk": newUserPk})
	}) + ";"
	return nil
}
