package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

func Init(conf *config.AppConfig) *gorm.DB {
	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
		conf.VAULTS_DB_USER,
		conf.VAULTS_DB_PASSWORD,
		conf.VAULTS_DB_HOST,
		conf.VAULTS_DB_PORT,
		conf.VAULTS_DB_NAME,
	)

	if db, err := gorm.Open(postgres.Open(dsn), &gormConfig); err != nil {
		panic(err)
	} else {
		return db
	}
}
