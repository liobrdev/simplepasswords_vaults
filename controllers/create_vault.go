package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type CreateVaultRequestBody struct {
	UserSlug   string `json:"user_slug"`
	VaultTitle string `json:"vault_title"`
}

func (H Handler) CreateVault(c *fiber.Ctx) error {
	body := CreateVaultRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.CreateVault, utils.ErrorParse, err.Error())
	}

	if !utils.SlugRegexp.MatchString(body.UserSlug) {
		return utils.RespondWithError(c, 400, utils.CreateVault, utils.ErrorUserSlug, body.UserSlug)
	}

	if body.VaultTitle == "" {
		return utils.RespondWithError(c, 400, utils.CreateVault, utils.ErrorVaultTitle, "")
	}

	if len(body.VaultTitle) > 255 {
		return utils.RespondWithError(c, 400, utils.CreateVault, utils.ErrorVaultTitle, "Too long")
	}

	var vault models.Vault

	if vaultSlug, err := utils.GenerateSlug(16); err != nil {
		return utils.RespondWithError(
			c, 500, utils.CreateVault,"Failed to generate `vault.Slug`.", err.Error(),
		)
	} else {
		vault.Slug = vaultSlug
	}

	vault.UserSlug = body.UserSlug
	vault.Title = body.VaultTitle

	if result := H.DB.Create(&vault); result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: vaults.title, vaults.user_slug" {
			return utils.RespondWithError(
				c, 409, utils.CreateVault, utils.ErrorDuplicateVault, result.Error.Error(),
			)
		}

		return utils.RespondWithError(
			c, 500, utils.CreateVault, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c, 500, utils.CreateVault, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(204)
}
