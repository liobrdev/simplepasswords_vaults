package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type UpdateVaultRequestBody struct {
	Title string `json:"vault_title"`
}

func (H Handler) UpdateVault(c *fiber.Ctx) error {
	body := UpdateVaultRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.UpdateVault, utils.ErrorParse, err.Error())
	}

	if body.Title == "" {
		return utils.RespondWithError(c, 400, utils.UpdateVault, utils.ErrorVaultTitle, "")
	}

	if len(body.Title) > 255 {
		return utils.RespondWithError(c, 400, utils.UpdateVault, utils.ErrorVaultTitle, "Too long")
	}

	slug := c.Params("slug")

	if result := H.DB.Model(&models.Vault{}).Where("slug = ?", slug).Update("title", body.Title);
	result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.UpdateVault, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c, 404, utils.UpdateVault, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c, 500, utils.UpdateVault, "result.RowsAffected > 1", strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(204)
}
