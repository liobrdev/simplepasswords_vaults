package secrets

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type CreateSecretRequestBody struct {
	UserSlug     string `json:"user_slug"`
	VaultSlug    string `json:"vault_slug"`
	EntrySlug    string `json:"entry_slug"`
	SecretLabel  string `json:"secret_label"`
	SecretString string `json:"secret_string"`
}

func (h handler) CreateSecret(c *fiber.Ctx) error {
	body := CreateSecretRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if !utils.SlugRegexp.MatchString(body.UserSlug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorUserSlug),
			body.UserSlug,
		)
	}

	if !utils.SlugRegexp.MatchString(body.VaultSlug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorVaultSlug),
			body.VaultSlug,
		)
	}

	if !utils.SlugRegexp.MatchString(body.EntrySlug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorEntrySlug),
			body.EntrySlug,
		)
	}

	if body.SecretLabel == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorSecretLabel),
			"",
		)
	}

	if len(body.SecretLabel) > 255 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorSecretLabel),
			fmt.Sprintf("Too long (%d > 255)", len(body.SecretLabel)),
		)
	}

	if body.SecretString == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorSecretString),
			"",
		)
	}

	if len(body.SecretString) > 1000 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateSecret,
			string(utils.ErrorSecretString),
			fmt.Sprintf("Too long (%d > 1000)", len(body.SecretString)),
		)
	}

	var secret models.Secret

	if secretSlug, err := utils.GenerateSlug(32); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusInternalServerError,
			utils.CreateSecret,
			"Failed to generate `secret.Slug`.",
			err.Error(),
		)
	} else {
		secret.Slug = secretSlug
	}

	secret.UserSlug = body.UserSlug
	secret.VaultSlug = body.VaultSlug
	secret.EntrySlug = body.EntrySlug
	secret.Label = body.SecretLabel
	secret.String = body.SecretString

	if result := h.DB.Create(&secret); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.CreateSecret,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
