package db

import (
	"strings"

	"github.com/atomicals-core/atomicals/DB/postsql"
)

func (m *Postgres) FtUTXOsByLocationID(locationID string) ([]*postsql.UTXOFtInfo, error) {
	// read from cache first, when this txID is in TxCache, this UTXONftInfo must in UserNftInfoCache; otherwise, this UTXONftInfo is not exist
	// entities, ok := m.UserFtInfoCache[locationID]
	// if ok {
	// 	return entities, nil
	// }

	var entity []*postsql.UTXOFtInfo
	dbTx := m.Where("location_id = ?", locationID).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	var res []*postsql.UTXOFtInfo
	for _, UTXO := range entity {
		res = append(res, &postsql.UTXOFtInfo{
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
	}
	return res, nil
}

func (m *Postgres) DistributedFtByName(tickerName string) (*postsql.GlobalDistributedFt, error) {
	var entity *postsql.GlobalDistributedFt
	dbTx := m.Where("ticker_name = ?", tickerName).Find(&entity)
	if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "record not found") {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, nil
	}
	return &postsql.GlobalDistributedFt{
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
	}, nil
}
