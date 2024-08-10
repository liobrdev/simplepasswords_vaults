package app

import (
	"fmt"
	"log"
	"net"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/database"
	"github.com/liobrdev/simplepasswords_vaults/routes"
)

func CreateApp(conf *config.AppConfig) (app *fiber.App, db *gorm.DB) {
	var proxyHeader string
	var trustedProxies []string

	if conf.BEHIND_PROXY {
		if err := parseIPs(&conf.PROXY_IP_ADDRESSES); err != nil {
			log.Fatal(err)
		} else {
			proxyHeader = "X-Forwarded-For"
			trustedProxies = conf.PROXY_IP_ADDRESSES
		}
	}

	app = fiber.New(fiber.Config{
		CaseSensitive:           true,
		EnableTrustedProxyCheck: conf.BEHIND_PROXY,
		ProxyHeader:             proxyHeader,
		TrustedProxies:          trustedProxies,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		Prefork:                 true,
	})

	db = database.Init(conf)
	routes.Register(app, db, conf)

	return app, db
}

func parseIPs(ipStrings *[]string) error {
	for _, s := range (*ipStrings) {
		if net.ParseIP(s) == nil {
			return fmt.Errorf("invalid IP: %s", s)
		}
	}

	return nil
}
