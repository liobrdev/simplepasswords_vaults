package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

func Init(conf *config.AppConfig) *gorm.DB {
	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	if db, err := gorm.Open(sqlite.Open("./test_db/vaults.sqlite"), &gormConfig); err != nil {
		panic(err)
	} else {
		return db
	}
}
