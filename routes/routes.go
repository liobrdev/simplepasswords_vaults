package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/controllers"
)

func Register(app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	H := controllers.Handler{DB: db, Conf: conf}
	app.Use(H.AuthorizeRequest)

	api := app.Group("/api")
	
	if H.Conf.ENVIRONMENT == "testing" {
		api.Get("/restricted", H.Restricted)
	}

	usersApi := api.Group("/users")
	usersApi.Post("/", H.CreateUser)
	usersApi.Get("/:slug", H.RetrieveUser)

	vaultsApi := api.Group("/vaults")
	vaultsApi.Post("/", H.CreateVault)
	vaultsApi.Get("/:slug", H.RetrieveVault)
	vaultsApi.Patch("/:slug", H.UpdateVault)
	// vaultsApi.Delete("/:slug", H.DeleteVault)
}
