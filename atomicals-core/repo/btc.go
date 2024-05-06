package repo

import (
	"strings"

	"github.com/atomicals-go/atomicals-core/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) InsertBtcTx(btcTx *postsql.BtcTx) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(&postsql.BtcTx{
			TxID:        btcTx.TxID,
			BlockHeight: btcTx.BlockHeight,
		})
	}) + ";"
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
		return -1, nil
	}
	return entity.BlockHeight, nil
}
