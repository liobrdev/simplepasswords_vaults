package controllers

import (
	"fmt"
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
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateVault,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if body.Title == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateVault,
			string(utils.ErrorVaultTitle),
			"",
		)
	}

	if len(body.Title) > 255 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateVault,
			string(utils.ErrorVaultTitle),
			fmt.Sprintf("Too long (%d > 255)", len(body.Title)),
		)
	}

	slug := c.Params("slug")

	if result := H.DB.Model(&models.Vault{}).Where("slug = ?", slug).Update(
		"title",
		body.Title,
	); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.UpdateVault,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c,
			fiber.StatusNotFound,
			utils.UpdateVault,
			string(utils.ErrorNoRowsAffected),
			"Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.UpdateVault,
			"result.RowsAffected > 1",
			strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
