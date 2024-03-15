package main

import (
	"github.com/atomicals-core/atomicals/DB/postsql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})

func InitModels(db *gorm.DB) {
	(&postsql.UserNftInfo{}).Init(db)
	(&postsql.GlobalDirectFt{}).Init(db)
	(&postsql.GlobalDistributedFt{}).Init(db)
	(&postsql.GlobalNftContainer{}).Init(db)
	(&postsql.GlobalNftRealm{}).Init(db)
	(&postsql.Location{}).Init(db)
	(&postsql.UserFtInfo{}).Init(db)
}

func AutoMigrate(db *gorm.DB) {
	(&postsql.UserNftInfo{}).AutoMigrate(db)
	(&postsql.GlobalDirectFt{}).AutoMigrate(db)
	(&postsql.GlobalDistributedFt{}).AutoMigrate(db)
	(&postsql.GlobalNftContainer{}).AutoMigrate(db)
	(&postsql.GlobalNftRealm{}).AutoMigrate(db)
	(&postsql.Location{}).AutoMigrate(db)
	(&postsql.UserFtInfo{}).AutoMigrate(db)
}

func main() {
	InitModels(DB)
	AutoMigrate(DB)
}
