package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	SQLRaw string
}

func (m *Postgres) UpdateCurrentHeightAndExecAllSql(height, txIndex int64) error {
	entity := &postsql.Location{}
	dbTx := m.Take(&entity)
	if dbTx.Error != nil && dbTx.Error != gorm.ErrRecordNotFound {
		return dbTx.Error
	}
	entity.Owner = "atomicals"
	entity.Height = height
	entity.TxIndex = txIndex
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(entity)
	})
	m.SQLRaw = m.SQLRaw + sql + ";"
	dbTx = m.Exec(m.SQLRaw)
	if dbTx.Error != nil {
		if !strings.Contains(dbTx.Error.Error(), "duplicate key value violates unique constraint") {
			return dbTx.Error
		}
	}
	m.SQLRaw = ""
	return nil
}

func (m *Postgres) CurrentLocation() (*postsql.Location, error) {
	entity := &postsql.Location{}
	dbTx := m.Find(&entity)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return entity, nil
}
