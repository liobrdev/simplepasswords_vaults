package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type CreateSecretRequestBody struct {
	UserSlug			 string `json:"user_slug"`
	VaultSlug			 string `json:"vault_slug"`
	EntrySlug			 string `json:"entry_slug"`
	SecretLabel		 string `json:"secret_label"`
	SecretString	 string `json:"secret_string"`
}

func (H Handler) CreateSecret(c *fiber.Ctx) error {
	body := CreateSecretRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorParse, err.Error())
	}

	if !utils.SlugRegexp.MatchString(body.UserSlug) {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorUserSlug, body.UserSlug)
	}

	if !utils.SlugRegexp.MatchString(body.VaultSlug) {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorVaultSlug, body.VaultSlug)
	}

	if !utils.SlugRegexp.MatchString(body.EntrySlug) {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorEntrySlug, body.EntrySlug)
	}

	if body.SecretLabel == "" {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorSecretLabel, "")
	}

	if len(body.SecretLabel) > 255 {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorSecretLabel, "Too long")
	}

	if body.SecretString == "" {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorSecretString, "")
	}

	if len(body.SecretString) > 1000 {
		return utils.RespondWithError(c, 400, utils.CreateSecret, utils.ErrorSecretString, "Too long")
	}

	var secret models.Secret

	if secretSlug, err := utils.GenerateSlug(32); err != nil {
		return utils.RespondWithError(
			c, 500, utils.CreateSecret, "Failed to generate `secret.Slug`.", err.Error(),
		)
	} else {
		secret.Slug = secretSlug
	}

	var count int64

	if result := H.DB.Model(&models.Secret{}).Where("entry_slug = ?", body.EntrySlug).Count(&count);
	result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.CreateSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else {
		secret.Priority = uint8(count)
	}

	secret.UserSlug = body.UserSlug
	secret.VaultSlug = body.VaultSlug
	secret.EntrySlug = body.EntrySlug
	secret.Label = body.SecretLabel
	secret.String = body.SecretString

	if result := H.DB.Create(&secret); result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.CreateSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	}

	return c.SendStatus(204)
}
