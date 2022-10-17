package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

func Init(conf *config.AppConfig) *gorm.DB {
	var db *gorm.DB

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	if conf.GO_FIBER_ENVIRONMENT != "production" {
		db = openDbSession("sqlite", "./test_db/vaults.sqlite", &gormConfig)
	} else {
		db = openDbSession("postgres", fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
			conf.GO_FIBER_VAULTS_DB_USER,
			conf.GO_FIBER_VAULTS_DB_PASSWORD,
			conf.GO_FIBER_VAULTS_DB_HOST,
			conf.GO_FIBER_VAULTS_DB_PORT,
			conf.GO_FIBER_VAULTS_DB_NAME,
		), &gormConfig)
	}

	return db
}

func openDbSession(driver string, dsn string, gormConfig *gorm.Config) (db *gorm.DB) {
	var err error

	if driver == "sqlite" {
		if db, err = gorm.Open(sqlite.Open(dsn), gormConfig); err != nil {
			panic(err)
		}
	} else if driver == "postgres" {
		if db, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
			panic(err)
		}
	} else {
		panic("Unsupported database driver: " + driver)
	}

	return
}
