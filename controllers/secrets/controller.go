package secrets

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

type handler struct {
	DB   *gorm.DB
	Conf *config.AppConfig
}

func RegisterSecrets(api *fiber.Router, db *gorm.DB, conf *config.AppConfig) {
	h := handler{db, conf}
	secretsApi := (*api).Group("/secrets")
	secretsApi.Post("/", h.CreateSecret)
	secretsApi.Patch("/:slug", h.UpdateSecret)
	secretsApi.Delete("/:slug", h.DeleteSecret)
}
