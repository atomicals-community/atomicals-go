package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertFtUTXO(UTXO *postsql.UTXOFtInfo) error {
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

func (m *Postgres) UpdateDistributedFtAmount(tickerName string, mintTimes int64) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", tickerName).Update("minted_times", mintTimes)
	}) + ";"
	return nil
}

func (m *Postgres) InsertDistributedFt(entity *postsql.GlobalDistributedFt) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(entity)
	}) + ";"
	return nil
}
