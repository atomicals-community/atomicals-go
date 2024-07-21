package main

import (
	"github.com/atomicals-go/pkg/conf"
	"github.com/atomicals-go/repo/postsql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitModels(db *gorm.DB) {
	(&postsql.GlobalDirectFt{}).Init(db)
	(&postsql.GlobalDistributedFt{}).Init(db)
	(&postsql.UTXOFtInfo{}).Init(db)
	(&postsql.UTXONftInfo{}).Init(db)
	(&postsql.ModInfo{}).Init(db)
	(&postsql.PaymentInfo{}).Init(db)
	(&postsql.AtomicalsTx{}).Init(db)
	(&postsql.Location{}).Init(db)
	(&postsql.BloomFilter{}).Init(db)
	(&postsql.StatisticTx{}).Init(db)
}

func AutoMigrate(db *gorm.DB) {
	(&postsql.GlobalDirectFt{}).AutoMigrate(db)
	(&postsql.GlobalDistributedFt{}).AutoMigrate(db)
	(&postsql.UTXOFtInfo{}).AutoMigrate(db)
	(&postsql.UTXONftInfo{}).AutoMigrate(db)
	(&postsql.ModInfo{}).AutoMigrate(db)
	(&postsql.PaymentInfo{}).AutoMigrate(db)
	(&postsql.AtomicalsTx{}).AutoMigrate(db)
	(&postsql.Location{}).AutoMigrate(db)
	(&postsql.BloomFilter{}).AutoMigrate(db)
	(&postsql.StatisticTx{}).AutoMigrate(db)
}

func main() {
	conf, err := conf.ReadJSONFromJSFile("./config.json")
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
