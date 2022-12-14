package main

import (
	"log"

	"github.com/liobrdev/simplepasswords_vaults/app"
	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/models"
)

func main() {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		log.Fatalln("Failed to load config from environment:", err)
	}

	app, db := app.CreateApp(&conf)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Vault{},
		&models.Entry{},
		&models.Secret{},
	); err != nil {
		log.Fatalln("Failed database auto-migrate:", err)
	}

	log.Fatal(app.Listen(conf.GO_FIBER_SERVER_HOST + ":" + conf.GO_FIBER_SERVER_PORT))
}
