package main

import (
	"github.com/atomicals-core/atomicals/DB/postsql"
	"github.com/atomicals-core/pkg/conf"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitModels(db *gorm.DB) {
	(&postsql.Location{}).Init(db)
	(&postsql.GlobalDirectFt{}).Init(db)
	(&postsql.GlobalDistributedFt{}).Init(db)
	(&postsql.UTXOFtInfo{}).Init(db)
	(&postsql.UTXONftInfo{}).Init(db)
}

func AutoMigrate(db *gorm.DB) {
	(&postsql.Location{}).AutoMigrate(db)
	(&postsql.GlobalDirectFt{}).AutoMigrate(db)
	(&postsql.GlobalDistributedFt{}).AutoMigrate(db)
	(&postsql.UTXOFtInfo{}).AutoMigrate(db)
	(&postsql.UTXONftInfo{}).AutoMigrate(db)
}

func main() {
	conf, err := conf.ReadJSONFromJSFile("../../../../conf/config.json")
	if err != nil {
		panic(err)
	}
	DB, err := gorm.Open(postgres.Open(conf.SqlDNS), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		panic(err)
	}
	InitModels(DB)
	// AutoMigrate(DB)
}
