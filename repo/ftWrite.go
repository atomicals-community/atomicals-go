package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertFtUTXO(UTXO *postsql.UTXOFtInfo) error {
	m.addFtLocationID(UTXO.LocationID)
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(UTXO)
	}) + ";"
	return nil
}

func (m *Postgres) DeleteFtUTXO(locationID string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.UTXOFtInfo{}).Where("location_id = ?", locationID).Delete(&postsql.UTXOFtInfo{})
	}) + ";"
	return nil
}

func (m *Postgres) UpdateDistributedFt(entity *postsql.GlobalDistributedFt) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", entity.TickerName).Updates(map[string]interface{}{"minted_times": entity.MintedTimes})
	}) + ";"
	return nil
}

func (m *Postgres) InsertDistributedFt(entity *postsql.GlobalDistributedFt) error {
	m.addDistributedFt(entity.TickerName)
	m.addFtLocationID(entity.LocationID)
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(entity)
	}) + ";"
	return nil
}

func (m *Postgres) InsertDirectFtUTXO(entity *postsql.GlobalDirectFt) error {
	m.addDirectFt(entity.TickerName)
	m.addFtLocationID(entity.LocationID)
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(entity)
	}) + ";"
	return nil
}
