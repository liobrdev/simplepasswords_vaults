package vaults

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

type handler struct {
	DB   *gorm.DB
	Conf *config.AppConfig
}

func RegisterVaults(api *fiber.Router, db *gorm.DB, conf *config.AppConfig) {
	h := handler{db, conf}
	vaultsApi := (*api).Group("/vaults")
	vaultsApi.Post("/", h.CreateVault)
	vaultsApi.Get("/:slug", h.RetrieveVault)
	vaultsApi.Patch("/:slug", h.UpdateVault)
	vaultsApi.Delete("/:slug", h.DeleteVault)
}
