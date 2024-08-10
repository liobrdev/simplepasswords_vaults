package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/controllers"
)

func Register(app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	H := controllers.Handler{DB: db, Conf: conf}
	api := app.Group("/api")

	usersApi := api.Group("/users")
	usersApi.Post("/", H.CreateUser)
	usersApi.Get("/:slug", H.RetrieveUser)

	// authApi := api.Group("/auth")
	// authApi.Post("/create_account", H.CreateAccount)
	// authApi.Post("/first_factor", H.AuthFirstFactor)
	// authApi.Post("/second_factor", H.AuthSecondFactor)

	// if H.Conf.ENVIRONMENT == "testing" {
	// 	authApi.Get("/restricted", H.AuthorizeRequest, H.Restricted)
	// }
}

// func Register(app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
// 	api := app.Group("/api")
// 	users.RegisterUsers(&api, db, conf)
// 	vaults.RegisterVaults(&api, db, conf)
// 	entries.RegisterEntries(&api, db, conf)
// 	secrets.RegisterSecrets(&api, db, conf)
// }
