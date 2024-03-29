package db

import (
	"github.com/atomicals-core/atomicals/DB/postsql"
	"github.com/atomicals-core/pkg/log"
	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
	SQLRaw string

	RealmCache     map[string]map[string]bool // key: realmName subrealmName
	ContainerCache map[string]map[string]bool // key: containerName itemNum
}

func (m *Postgres) initDBCache() error {
	log.Log.Infof("init DB Cache...")
	entity := make([]*postsql.UTXONftInfo, 0)
	dbTx := m.Where("realm_name != ?", "").Find(&entity)
	if dbTx.Error != nil && dbTx.Error != gorm.ErrRecordNotFound {
		return dbTx.Error
	}
	for _, v := range entity {
		if _, exist := m.RealmCache[v.RealmName]; exist {
			m.RealmCache[v.RealmName][v.SubRealmName] = true
		} else {
			m.RealmCache[v.RealmName] = make(map[string]bool)
			m.RealmCache[v.RealmName][v.SubRealmName] = true
		}
	}
	dbTx = m.Where("container_name != ?", "").Find(&entity)
	if dbTx.Error != nil && dbTx.Error != gorm.ErrRecordNotFound {
		return dbTx.Error
	}
	for _, v := range entity {
		if _, exist := m.ContainerCache[v.ContainerName]; exist {
			m.ContainerCache[v.ContainerName][v.Dmitem] = true
		} else {
			m.ContainerCache[v.ContainerName] = make(map[string]bool)
			m.ContainerCache[v.ContainerName][v.Dmitem] = true
		}
	}
	return nil
}

func (m *Postgres) UpdateCurrentHeightAndExecAllSql(height int64, txID string) error {
	entity := &postsql.Location{}
	dbTx := m.Take(&entity)
	if dbTx.Error != nil && dbTx.Error != gorm.ErrRecordNotFound {
		return dbTx.Error
	}
	entity.Owner = "atomicals"
	entity.Height = height
	entity.TxID = txID
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Save(entity)
	})
	m.SQLRaw = m.SQLRaw + sql + ";"
	dbTx = m.Exec(m.SQLRaw)
	if dbTx.Error != nil {
		return dbTx.Error
	}
	m.SQLRaw = ""
	return nil
}

func (m *Postgres) CurrentHeitht() (int64, error) {
	entity := &postsql.Location{}
	dbTx := m.Find(&entity)
	if dbTx.Error != nil {
		return -1, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return -1, gorm.ErrRecordNotFound
	}
	return entity.Height, nil
}
