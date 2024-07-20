package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

func (m *Postgres) Location() (*postsql.Location, error) {
	entity := &postsql.Location{}
	dbTx := m.Order("id desc").First(&entity)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return entity, nil
}

func (m *Postgres) AtomicalsTx(txID string) (*postsql.AtomicalsTx, error) {
	var entity *postsql.AtomicalsTx
	dbTx := m.Where("tx_id = ?", txID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (m *Postgres) AtomicalsTxHeight(txID string) (int64, error) {
	var entity *postsql.AtomicalsTx
	dbTx := m.Where("tx_id = ?", txID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return -1, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return -1, gorm.ErrRecordNotFound
	}
	return entity.BlockHeight, nil
}
