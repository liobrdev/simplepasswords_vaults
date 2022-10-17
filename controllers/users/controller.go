package users

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

type handler struct {
	DB   *gorm.DB
	Conf *config.AppConfig
}

func RegisterUsers(api *fiber.Router, db *gorm.DB, conf *config.AppConfig) {
	h := handler{db, conf}
	usersApi := (*api).Group("/users")
	usersApi.Post("/", h.CreateUser)
	usersApi.Get("/:slug", h.RetrieveUser)
}
