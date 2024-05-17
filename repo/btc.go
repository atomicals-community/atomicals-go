package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertBtcTx(btcTx *postsql.BtcTx) error {
	dbTx := m.Save(&postsql.BtcTx{
		TxID:        btcTx.TxID,
		BlockHeight: btcTx.BlockHeight,
	})
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "duplicate key value violates unique constraint") {
		return dbTx.Error
	}
	return nil
}

func (m *Postgres) BtcTx(txID string) (*postsql.BtcTx, error) {
	var entity *postsql.BtcTx
	dbTx := m.Where("tx_id = ?", txID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (m *Postgres) BtcTxHeight(txID string) (int64, error) {
	var entity *postsql.BtcTx
	dbTx := m.Where("tx_id = ?", txID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return -1, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return -1, gorm.ErrRecordNotFound
	}
	return entity.BlockHeight, nil
}

func (m *Postgres) InsertMod(mod *postsql.ModInfo) error {
	dbTx := m.Save(mod)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "duplicate key value violates unique constraint") {
		return dbTx.Error
	}
	return nil
}

func (m *Postgres) Mod(atomicalsID string) (*postsql.ModInfo, error) {
	var entity *postsql.ModInfo
	dbTx := m.Where("atomicals_id = ?", atomicalsID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}
