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
	
	vaultsApi := api.Group("/vaults")
	vaultsApi.Post("/", H.CreateVault)
	vaultsApi.Get("/", H.ListVaults)
	vaultsApi.Get("/:slug", H.RetrieveVault)
	vaultsApi.Patch("/:slug", H.UpdateVault)
	vaultsApi.Delete("/:slug", H.DeleteVault)

	entriesApi := api.Group("/entries")
	entriesApi.Post("/", H.CreateEntry)
	entriesApi.Get("/:slug", H.RetrieveEntry)
	entriesApi.Patch("/:slug", H.UpdateEntry)
	entriesApi.Delete("/:slug", H.DeleteEntry)

	secretsApi := api.Group("/secrets")
	secretsApi.Post("/", H.CreateSecret)
	secretsApi.Patch("/:slug", H.UpdateSecret, H.MoveSecret)
	secretsApi.Delete("/:slug", H.DeleteSecret)
}
