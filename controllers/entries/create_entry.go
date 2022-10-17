package entries

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type reqBodySecret struct {
	Label  string `json:"secret_label"`
	String string `json:"secret_string"`
}

type CreateEntryRequestBody struct {
	UserSlug   string          `json:"user_slug"`
	VaultSlug  string          `json:"vault_slug"`
	EntryTitle string          `json:"entry_title"`
	Secrets    []reqBodySecret `json:"secrets"`
}

func (h handler) CreateEntry(c *fiber.Ctx) error {
	body := CreateEntryRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateEntry,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if !utils.SlugRegexp.MatchString(body.UserSlug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateEntry,
			string(utils.ErrorUserSlug),
			body.UserSlug,
		)
	}

	if !utils.SlugRegexp.MatchString(body.VaultSlug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateEntry,
			string(utils.ErrorVaultSlug),
			body.VaultSlug,
		)
	}

	if body.EntryTitle == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateEntry,
			string(utils.ErrorEntryTitle),
			"",
		)
	}

	if len(body.EntryTitle) > 255 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateEntry,
			string(utils.ErrorEntryTitle),
			fmt.Sprintf("Too long (%d > 255)", len(body.EntryTitle)),
		)
	}

	if body.Secrets == nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateEntry,
			string(utils.ErrorSecrets),
			"",
		)
	}

	labels := map[string]bool{}

	for i, secretsLen := 0, len(body.Secrets); i < secretsLen; i++ {
		if body.Secrets[i].Label == "" || len(body.Secrets[i].Label) > 255 {
			return utils.RespondWithError(
				c,
				fiber.StatusBadRequest,
				utils.CreateEntry,
				string(utils.ErrorItemSecrets),
				fmt.Sprintf("secrets[%d].Label; len(secrets) == %d", i, secretsLen),
			)
		}

		if body.Secrets[i].String == "" || len(body.Secrets[i].String) > 1000 {
			return utils.RespondWithError(
				c,
				fiber.StatusBadRequest,
				utils.CreateEntry,
				string(utils.ErrorItemSecrets),
				fmt.Sprintf("secrets[%d].String; len(secrets) == %d", i, secretsLen),
			)
		}

		if _, ok := labels[body.Secrets[i].Label]; ok {
			return utils.RespondWithError(
				c,
				fiber.StatusBadRequest,
				utils.CreateEntry,
				string(utils.ErrorDuplicateSecrets),
				body.Secrets[i].Label,
			)
		}

		labels[body.Secrets[i].Label] = true
	}

	var entry models.Entry

	if entrySlug, err := utils.GenerateSlug(32); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusInternalServerError,
			utils.CreateEntry,
			"Failed to generate `entry.Slug`.",
			err.Error(),
		)
	} else {
		entry.Slug = entrySlug
	}

	entry.UserSlug = body.UserSlug
	entry.VaultSlug = body.VaultSlug
	entry.Title = body.EntryTitle

	if err := h.DB.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&entry); result.Error != nil {
			return result.Error
		}

		for i, secretsLen, secretSlug := 0, len(body.Secrets), ""; i < secretsLen; i++ {
			if slug, err := utils.GenerateSlug(32); err != nil {
				return fmt.Errorf("`secret.Slug` generation failed: %s", err.Error())
			} else {
				secretSlug = slug
			}

			if result := tx.Create(&models.Secret{
				Slug:      secretSlug,
				Label:     body.Secrets[i].Label,
				String:    body.Secrets[i].String,
				EntrySlug: entry.Slug,
				VaultSlug: entry.VaultSlug,
				UserSlug:  entry.UserSlug,
			}); result.Error != nil {
				return result.Error
			}
		}

		return nil
	}); err != nil {
		if errText := err.Error(); utils.FailedSecretSlugRegexp.MatchString(errText) {
			return utils.RespondWithError(
				c,
				fiber.StatusInternalServerError,
				utils.CreateEntry,
				errText,
				"",
			)
		}

		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.CreateEntry,
			string(utils.ErrorFailedDB),
			err.Error(),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
