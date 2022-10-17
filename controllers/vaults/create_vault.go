package vaults

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type CreateVaultRequestBody struct {
	UserSlug   string `json:"user_slug"`
	VaultTitle string `json:"vault_title"`
}

func (h handler) CreateVault(c *fiber.Ctx) error {
	body := CreateVaultRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateVault,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if !utils.SlugRegexp.MatchString(body.UserSlug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateVault,
			string(utils.ErrorUserSlug),
			body.UserSlug,
		)
	}

	if body.VaultTitle == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateVault,
			string(utils.ErrorVaultTitle),
			"",
		)
	}

	if len(body.VaultTitle) > 255 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateVault,
			string(utils.ErrorVaultTitle),
			fmt.Sprintf("Too long (%d > 255)", len(body.VaultTitle)),
		)
	}

	var vault models.Vault

	if vaultSlug, err := utils.GenerateSlug(32); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusInternalServerError,
			utils.CreateVault,
			"Failed to generate `vault.Slug`.",
			err.Error(),
		)
	} else {
		vault.Slug = vaultSlug
	}

	vault.UserSlug = body.UserSlug
	vault.Title = body.VaultTitle

	if result := h.DB.Create(&vault); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.CreateVault,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
