package repo

import (
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertFtUTXO(UTXO *postsql.UTXOFtInfo) error {
	// update TXCache first
	// _, ok := m.UserFtInfoCache[UTXO.LocationID]
	// if !ok {
	// 	m.UserFtInfoCache[UTXO.LocationID] = make([]*UTXOFtInfo, 0)
	// }
	// m.UserFtInfoCache[UTXO.LocationID] = append(m.UserFtInfoCache[UTXO.LocationID], UTXO)

	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.UTXOFtInfo{
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
	}) + ";"
	return nil
}

func (m *Postgres) DeleteFtUTXO(locationID string) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.UTXOFtInfo{}).Where("loaction_id = ?", locationID).Delete(&postsql.UTXOFtInfo{})
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
