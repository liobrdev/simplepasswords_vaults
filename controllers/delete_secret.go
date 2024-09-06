package controllers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) DeleteSecret(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.DeleteSecret, utils.ErrorSecretSlug, slug)
	}

	var secret models.Secret

	if result := H.DB.First(&secret, "slug = ?", slug); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.RespondWithError(c, 404, utils.DeleteSecret, utils.ErrorNotFound, slug)
		}

		return utils.RespondWithError(
			c, 500, utils.DeleteSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	}

	var secrets []string

	if result := H.DB.Raw(
		"SELECT slug FROM secrets WHERE entry_slug = ? AND priority > ?",
		secret.EntrySlug, secret.Priority,
	).Scan(&secrets); result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.DeleteSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	}

	if err := H.DB.Transaction(func(tx *gorm.DB) error {
		if result := tx.Exec(
			"UPDATE secrets SET priority = ?, updated_at = ? WHERE slug IN ?",
			gorm.Expr("priority - 1"), time.Now().UTC(), secrets,
		); result.Error != nil {
			return result.Error
		}

		if result := tx.Delete(&secret); result.Error != nil {
			return result.Error
		}

		return nil
	}); err != nil {
		return utils.RespondWithError(c, 500, utils.MoveSecret, utils.ErrorFailedDB, err.Error())
	}

	return c.SendStatus(204)
}
