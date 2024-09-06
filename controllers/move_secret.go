package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) MoveSecret(c *fiber.Ctx) error {
	body := UpdateSecretRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.MoveSecret, utils.ErrorParse, err.Error())
	}

	if body.Priority == "" {
		return utils.RespondWithError(c, 400, utils.MoveSecret, utils.ErrorSecretPriority, "")
	}

	var newPriority int
	var err error

	if newPriority, err = strconv.Atoi(body.Priority); err != nil {
		return utils.RespondWithError(c, 400, utils.MoveSecret, utils.ErrorSecretPriority, err.Error())
	}

	if !utils.SlugRegexp.MatchString(body.EntrySlug) {
		return utils.RespondWithError(c, 400, utils.MoveSecret, utils.ErrorEntrySlug, body.EntrySlug)
	}

	var secrets []models.Secret
	slug := c.Params("slug")

	if result := H.DB.Where("entry_slug = ?", body.EntrySlug).Find(&secrets); result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.MoveSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if result.RowsAffected == 0 {
		return utils.RespondWithError(c, 404, utils.MoveSecret, utils.ErrorNotFound, "No secrets found")
	} else if result.RowsAffected == 1 {
		if result := H.DB.Model(models.Secret{}).Where("slug = ?", slug).Update("priority", 0);
		result.Error != nil {
			return utils.RespondWithError(
				c, 500, utils.MoveSecret, utils.ErrorFailedDB, result.Error.Error(),
			)
		}

		return c.SendStatus(204)
	}

	if newPriority > len(secrets) {
		newPriority = len(secrets) - 1
	} else if newPriority < 0 {
		newPriority = 0
	}

	var thisSecret *models.Secret

	for _, secret := range secrets {
		if secret.Slug == slug {
			thisSecret = &secret
			break
		}
	}

	if thisSecret == nil {
		return utils.RespondWithError(c, 404, utils.MoveSecret, utils.ErrorSecretSlug, "Not found")
	}

	var expr clause.Expr
	secretsToUpdate := []string{}

	if newPriority < int(thisSecret.Priority) {
		expr = gorm.Expr("priority + 1")

		for _, secret := range secrets {
			if secret.Priority < thisSecret.Priority && secret.Priority >= uint8(newPriority) {
				secretsToUpdate = append(secretsToUpdate, secret.Slug)
			}
		}
	} else if newPriority > int(thisSecret.Priority) {
		expr = gorm.Expr("priority - 1")

		for _, secret := range secrets {
			if secret.Priority > thisSecret.Priority && secret.Priority <= uint8(newPriority) {
				secretsToUpdate = append(secretsToUpdate, secret.Slug)
			}
		}
	} else {
		return c.SendStatus(204)
	}

	if err := H.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()

		if result := tx.Exec(
			"UPDATE secrets SET priority = ?, updated_at = ? WHERE slug IN ?",
			expr, now, secretsToUpdate,
		); result.Error != nil {
			return result.Error
		}

		if result := tx.Exec(
			"UPDATE secrets SET priority = ?, updated_at = ? WHERE slug = ?",
			newPriority, now, thisSecret.Slug,
		); result.Error != nil {
			return result.Error
		}

		return nil
	}); err != nil {
		return utils.RespondWithError(c, 500, utils.MoveSecret, utils.ErrorFailedDB, err.Error())
	}

	return c.SendStatus(204)
}
