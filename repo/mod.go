package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) ModHistory(atomicalsID string, height int64) ([]*postsql.ModInfo, error) {
	var entities []*postsql.ModInfo
	dbTx := m.Where("atomicals_id = ? and height >= ?", atomicalsID, height).Order("id").Find(&entities)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entities, nil
}
