package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type ListVaultsResponseBody struct {
	Vaults []models.Vault `json:"vaults"`
}

func (H Handler) ListVaults(c *fiber.Ctx) error {
	slug := c.Get("User-Slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.ListVaults, utils.ErrorUserSlug, slug)
	}

	var vaults []models.Vault

	if result := H.DB.Find(&vaults, "user_slug = ?", slug); result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.ListVaults, utils.ErrorFailedDB, result.Error.Error(),
		)
	}

	return c.Status(200).JSON(&ListVaultsResponseBody{ Vaults: vaults })
}
