package main

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"

	"github.com/liobrdev/simplepasswords_vaults/app"
	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/database"
	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/routes"
)

func main() {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		log.Fatalln("Failed to load config from environment:", err)
	}

	app := app.CreateApp(&conf)
	db := database.Init(&conf)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Vault{},
		&models.Entry{},
		&models.Secret{},
	); err != nil {
		log.Fatalln("Failed database auto-migrate:", err)
	}

	app.Use(healthcheck.New())
	routes.Register(app, db, &conf)

	log.Fatal(app.Listen(conf.VAULTS_HOST + ":" + conf.VAULTS_PORT))
}
