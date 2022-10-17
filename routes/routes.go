package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/controllers/entries"
	"github.com/liobrdev/simplepasswords_vaults/controllers/secrets"
	"github.com/liobrdev/simplepasswords_vaults/controllers/users"
	"github.com/liobrdev/simplepasswords_vaults/controllers/vaults"
)

func RegisterAPI(app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	api := app.Group("/api")
	users.RegisterUsers(&api, db, conf)
	vaults.RegisterVaults(&api, db, conf)
	entries.RegisterEntries(&api, db, conf)
	secrets.RegisterSecrets(&api, db, conf)
}
