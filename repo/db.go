package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	SQLRaw string

	bloomFilter map[string]*bloomFilterInfo
}

func (m *Postgres) ExecAllSql(blockHeight, txIndex int64, txID, operation string) error {
	m.InsertBtcTx(&postsql.BtcTx{
		BlockHeight: blockHeight,
		TxIndex:     txIndex,
		TxID:        txID,
		Operation:   operation,
		Description: m.SQLRaw,
	})

	for name, v := range m.bloomFilter {
		if v.needUpdate {
			m.UpdateBloomFilter(name, v.filter)
		}
	}

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.Location{}).Where("name = ?", "atomicals").Updates(map[string]interface{}{"block_height": blockHeight, "tx_index": txIndex})
	})
	m.SQLRaw = m.SQLRaw + sql + ";"
	dbTx := m.Exec(m.SQLRaw)
	m.SQLRaw = ""
	if dbTx.Error != nil {
		return dbTx.Error
	}
	for _, v := range m.bloomFilter {
		v.needUpdate = false
	}
	return nil
}

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

// if dbTx.Error != nil && !strings.Contains(dbTx.Error.Error(), "duplicate key value violates unique constraint") {
func (m *Postgres) InsertBtcTx(btcTx *postsql.BtcTx) error {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(btcTx)
	})
	m.SQLRaw = m.SQLRaw + sql + ";"
	return nil
}

func (m *Postgres) DeleteBtcTxUntil(blockHeight int64) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.BtcTx{}).Unscoped().Where("block_height = ? and operation = ?", blockHeight, "").Delete(&postsql.UTXOFtInfo{})
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
		return -1, gorm.ErrRecordNotFound
	}
	return entity.BlockHeight, nil
}

func (m *Postgres) InsertMod(mod *postsql.ModInfo) error {
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return m.Save(mod)
	}) + ";"
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
