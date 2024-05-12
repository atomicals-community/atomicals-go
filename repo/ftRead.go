package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) FtUTXOsByLocationID(locationID string) ([]*postsql.UTXOFtInfo, error) {
	var entity []*postsql.UTXOFtInfo
	dbTx := m.Model(postsql.UTXOFtInfo{}).Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}

func (m *Postgres) DistributedFtByName(tickerName string) (*postsql.GlobalDistributedFt, error) {
	var entity *postsql.GlobalDistributedFt
	dbTx := m.Model(postsql.GlobalDistributedFt{}).Where("ticker_name = ?", tickerName).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}
