package app

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/config"
)

func CreateApp(conf *config.AppConfig) (app *fiber.App) {
	return fiber.New(fiber.Config{
		CaseSensitive:           true,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		Prefork:                 true,
	})
}
