package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) AuthorizeRequest(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// Check Authorization header format
	if !utils.AuthHeaderRegexp.Match([]byte(authHeader)) {
		return utils.RespondWithError(c, 401, c.Get("Client-Operation"), utils.ErrorToken, authHeader)
	}

	authToken := authHeader[6:]

	if authToken != H.Conf.VAULTS_ACCESS_TOKEN {
		return utils.RespondWithError(c, 401, c.Get("Client-Operation"), utils.ErrorToken, authToken)
	}

	return c.Next()
}
