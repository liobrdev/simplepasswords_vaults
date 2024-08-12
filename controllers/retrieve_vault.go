package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) RetrieveVault(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.RetrieveVault, utils.ErrorVaultSlug, slug)
	}

	var vault models.Vault

	if result := H.DB.Preload("Entries").First(&vault, "slug = ?", slug); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.RespondWithError(c, 404, utils.RetrieveVault, utils.ErrorNotFound, slug)
		}

		return utils.RespondWithError(
			c, 500, utils.RetrieveVault, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c, 500, utils.RetrieveVault, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
		)
	}

	return c.Status(400).JSON(&vault)
}
