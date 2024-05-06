package postsql

import "gorm.io/gorm"

type DefaultModelImpl struct {
	Table string
	DB    *gorm.DB
}

func newDefaultModel(table string, db *gorm.DB) *DefaultModelImpl {
	return &DefaultModelImpl{
		Table: table,
		DB:    db,
	}
}

func (m *DefaultModelImpl) DropTable() error {
	return m.DB.Migrator().DropTable(m.Table)
}
func (m *DefaultModelImpl) CreateTable(model interface{}) error {
	return m.DB.AutoMigrate(model)
}
