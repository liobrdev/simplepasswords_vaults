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

func CreateApp(conf *config.AppConfig) (*fiber.App, *gorm.DB) {
	var proxyHeader string
	var trustedProxies []string

	if conf.GO_FIBER_BEHIND_PROXY {
		if err := parseIPs(&conf.GO_FIBER_PROXY_IP_ADDRESSES); err != nil {
			log.Fatal(err)
		} else {
			proxyHeader = "X-Forward-For"
			trustedProxies = conf.GO_FIBER_PROXY_IP_ADDRESSES
		}
	}

	app := fiber.New(fiber.Config{
		CaseSensitive:           true,
		EnableTrustedProxyCheck: conf.GO_FIBER_BEHIND_PROXY,
		ProxyHeader:             proxyHeader,
		TrustedProxies:          trustedProxies,
		JSONEncoder:             json.Marshal,
		JSONDecoder:             json.Unmarshal,
		Prefork:                 true,
	})

	db := database.Init(conf)
	routes.RegisterAPI(app, db, conf)

	return app, db
}

func parseIPs(ipStrings *[]string) error {
	for i, n := 0, len(*ipStrings); i < n; i++ {
		if net.ParseIP((*ipStrings)[i]) == nil {
			return fmt.Errorf("invalid IP: %s", (*ipStrings)[i])
		}
	}

	return nil
}
