package repo

import (
	"strings"

	"github.com/atomicals-go/repo/postsql"
	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	SQLRaw string
}

func (m *Postgres) ExecAllSql(location *postsql.Location) error {
	m.InsertBtcTx(&postsql.BtcTx{
		BlockHeight: location.BlockHeight,
		TxIndex:     location.TxIndex,
		TxID:        location.Txid,
	})

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.Location{}).Where("name = ?", "atomicals").Updates(map[string]interface{}{"block_height": location.BlockHeight, "tx_index": location.TxIndex})
	})
	m.SQLRaw = m.SQLRaw + sql + ";"
	dbTx := m.Exec(m.SQLRaw)
	m.SQLRaw = ""
	if dbTx.Error != nil {
		return dbTx.Error
	}
	return nil
}

func (m *Postgres) CurrentLocation() (*postsql.Location, error) {
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

func (m *Postgres) InsertBloomFilter(name string, filter *bloom.BloomFilter) error {
	data, err := filter.MarshalJSON()
	if err != nil {
		return err
	}
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return m.Save(&postsql.BloomFilter{
			Name: name,
			Data: data,
		})
	}) + ";"
	return nil
}

func (m *Postgres) UpdateBloomFilter(name string, filter *bloom.BloomFilter) error {
	data, err := filter.MarshalJSON()
	if err != nil {
		return err
	}
	m.SQLRaw += m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(postsql.BloomFilter{}).Where("name = ?", name).Update("data", data)
	}) + ";"
	return nil
}

func (m *Postgres) BloomFilter() (map[string]*bloom.BloomFilter, error) {
	entities := make([]*postsql.BloomFilter, 0)
	dbTx := m.Order("id desc").Find(&entities)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	filterMap := make(map[string]*bloom.BloomFilter, 0)
	for _, v := range entities {
		filter := bloom.NewWithEstimates(10000, 0.01)
		filter.UnmarshalJSON(v.Data)
		filterMap[v.Name] = filter
	}
	return filterMap, nil
}
