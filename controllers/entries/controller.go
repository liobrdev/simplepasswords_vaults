package entries

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

type handler struct {
	DB   *gorm.DB
	Conf *config.AppConfig
}

func RegisterEntries(api *fiber.Router, db *gorm.DB, conf *config.AppConfig) {
	h := handler{db, conf}
	entriesApi := (*api).Group("/entries")
	entriesApi.Post("/", h.CreateEntry)
	entriesApi.Get("/:slug", h.RetrieveEntry)
	entriesApi.Patch("/:slug", h.UpdateEntry)
	entriesApi.Delete("/:slug", h.DeleteEntry)
}
