package main

import (
	"github.com/atomicals-core/atomicals/DB/postsql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB, _ = gorm.Open(postgres.Open("host=127.0.0.1 user=postgres password=ZecreyProtocolDB@123 dbname=atomicals port=5432 sslmode=disable"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

func InitModels(db *gorm.DB) {
	(&postsql.Location{}).Init(db)
	(&postsql.GlobalNft{}).Init(db)
	(&postsql.GlobalDirectFt{}).Init(db)
	(&postsql.GlobalDistributedFt{}).Init(db)
	(&postsql.UserFtInfo{}).Init(db)
	(&postsql.UserNftInfo{}).Init(db)
}

func AutoMigrate(db *gorm.DB) {
	(&postsql.Location{}).AutoMigrate(db)
	(&postsql.GlobalNft{}).AutoMigrate(db)
	(&postsql.GlobalDirectFt{}).AutoMigrate(db)
	(&postsql.GlobalDistributedFt{}).AutoMigrate(db)
	(&postsql.UserFtInfo{}).AutoMigrate(db)
	(&postsql.UserNftInfo{}).AutoMigrate(db)
}

func main() {
	InitModels(DB)
	AutoMigrate(DB)
}
