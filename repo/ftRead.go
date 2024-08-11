package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
)

func (m *Postgres) FtUTXOsByUserPK(UserPK string) ([]*postsql.UTXOFtInfo, error) {
	var entity []*postsql.UTXOFtInfo
	dbTx := m.Model(postsql.UTXOFtInfo{}).Where("user_pk = ?", UserPK).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	// if dbTx.RowsAffected == 0 {
	// 	return nil, nil
	// }
	return entity, nil
}

func (m *Postgres) FtUTXOsByLocationID(locationID string) ([]*postsql.UTXOFtInfo, error) {
	if !m.testFtLocationID(locationID) {
		return nil, nil
	}
	var entity []*postsql.UTXOFtInfo
	dbTx := m.Model(postsql.UTXOFtInfo{}).Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}

func (m *Postgres) DistributedFtByName(tickerName string) (*postsql.GlobalDistributedFt, error) {
	// if !m.testDistributedFt(tickerName) {
	// 	return nil, nil
	// }
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

func (m *Postgres) DirectFtByName(tickerName string) (*postsql.GlobalDirectFt, error) {
	var entity *postsql.GlobalDirectFt
	dbTx := m.Model(postsql.GlobalDirectFt{}).Where("ticker_name = ?", tickerName).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return entity, nil
}

func (m *Postgres) FtUTXOsByID(offset, limit int) ([]*postsql.UTXOFtInfo, error) {
	var entity []*postsql.UTXOFtInfo
	dbTx := m.Model(postsql.UTXOFtInfo{}).Order("id").Offset(offset).Limit(limit).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}

func (m *Postgres) DistributedFtByID(offset, limit int) ([]*postsql.GlobalDistributedFt, error) {
	var entity []*postsql.GlobalDistributedFt
	dbTx := m.Model(postsql.GlobalDistributedFt{}).Order("id").Offset(offset).Limit(limit).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	return entity, nil
}
