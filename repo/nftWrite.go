package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertNftUTXO(UTXO *postsql.UTXONftInfo) error {
	if UTXO.RealmName != "" {
		m.addRealm(UTXO.RealmName)
	} else if UTXO.ContainerName != "" {
		m.addContainer(UTXO.ContainerName)
	}
	m.addNftLocationID(UTXO.LocationID)
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(UTXO)
	}) + ";"
	return nil
}

func (m *Postgres) TransferNftUTXO(oldLocationID, newLocationID, newUserPk string) error {
	m.addNftLocationID(newLocationID)
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.UTXONftInfo{}).Where("location_id = ?", oldLocationID).Updates(map[string]interface{}{"location_id": newLocationID, "user_pk": newUserPk})
	}) + ";"
	return nil
}
