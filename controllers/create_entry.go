package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type reqBodySecret struct {
	Label  	 string	`json:"secret_label"`
	String 	 string	`json:"secret_string"`
	Priority uint8 	`json:"secret_priority"`
}

type CreateEntryRequestBody struct {
	UserSlug   string          `json:"user_slug"`
	VaultSlug  string          `json:"vault_slug"`
	EntryTitle string          `json:"entry_title"`
	Secrets    []reqBodySecret `json:"secrets"`
}

func (H Handler) CreateEntry(c *fiber.Ctx) error {
	body := CreateEntryRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.CreateEntry, utils.ErrorParse, err.Error())
	}

	if !utils.SlugRegexp.MatchString(body.UserSlug) {
		return utils.RespondWithError(c, 400, utils.CreateEntry, utils.ErrorUserSlug, body.UserSlug)
	}

	if !utils.SlugRegexp.MatchString(body.VaultSlug) {
		return utils.RespondWithError(c, 400, utils.CreateEntry, utils.ErrorVaultSlug, body.VaultSlug)
	}

	if body.EntryTitle == "" {
		return utils.RespondWithError(c, 400, utils.CreateEntry, utils.ErrorEntryTitle, "")
	}

	if len(body.EntryTitle) > 255 {
		return utils.RespondWithError(c, 400, utils.CreateEntry, utils.ErrorEntryTitle, "Too long")
	}

	if body.Secrets == nil {
		return utils.RespondWithError(c, 400, utils.CreateEntry, utils.ErrorSecrets, "")
	}

	secretsLen := len(body.Secrets)
	labels := map[string]bool{}
	priorities := map[uint8]bool{}

	for i, secret := range(body.Secrets) {
		if secret.Label == "" || len(secret.Label) > 255 {
			return utils.RespondWithError(
				c, 400, utils.CreateEntry, utils.ErrorItemSecrets,
				fmt.Sprintf("secrets[%d].Label; len(secrets) == %d", i, secretsLen),
			)
		}

		if secret.String == "" || len(secret.String) > 1000 {
			return utils.RespondWithError(
				c, 400, utils.CreateEntry, utils.ErrorItemSecrets,
				fmt.Sprintf("secrets[%d].String; len(secrets) == %d", i, secretsLen),
			)
		}

		if _, ok := labels[secret.Label]; ok {
			return utils.RespondWithError(
				c, 400, utils.CreateEntry, utils.ErrorDuplicateSecretsLabel, secret.Label,
			)
		}

		if _, ok := priorities[secret.Priority]; ok {
			return utils.RespondWithError(
				c, 400, utils.CreateEntry, utils.ErrorDuplicateSecretsPriority, "",
			)
		}

		labels[secret.Label] = true
		priorities[secret.Priority] = true
	}

	var entry models.Entry

	if entrySlug, err := utils.GenerateSlug(16); err != nil {
		return utils.RespondWithError(
			c, 500, utils.CreateEntry, "Failed to generate `entry.Slug`.", err.Error(),
		)
	} else {
		entry.Slug = entrySlug
	}

	entry.UserSlug = body.UserSlug
	entry.VaultSlug = body.VaultSlug
	entry.Title = body.EntryTitle

	if err := H.DB.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&entry); result.Error != nil {
			return result.Error
		}

		password := c.Get(H.Conf.PASSWORD_HEADER_KEY)

		for _, secret := range(body.Secrets) {
			if slug, err := utils.GenerateSlug(16); err != nil {
				return fmt.Errorf("`secret.Slug` generation failed: %s", err.Error())
			} else if encryptedString, err := utils.Encrypt(secret.String, password); err != nil {
				return fmt.Errorf("`secret.String` encryption failed: %s", err.Error())
			} else if result := tx.Create(&models.Secret{
				Slug:      slug,
				Label:     secret.Label,
				String:    encryptedString,
				Priority:	 secret.Priority,
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
			return utils.RespondWithError(c, 500, utils.CreateEntry, errText, "")
		}

		return utils.RespondWithError(c, 500, utils.CreateEntry, utils.ErrorFailedDB, err.Error())
	}

	return c.SendStatus(204)
}
