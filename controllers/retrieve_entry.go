package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) RetrieveEntry(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.RetrieveEntry, utils.ErrorEntrySlug, slug)
	}

	var entry models.Entry

	if result := H.DB.Preload("Secrets", func(db *gorm.DB) *gorm.DB {
		return db.Order("secrets.priority, secrets.created_at")
	}).First(&entry, "slug = ?", slug); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.RespondWithError(c, 404, utils.RetrieveEntry, utils.ErrorNotFound, slug)
		}

		return utils.RespondWithError(
			c, 500, utils.RetrieveEntry, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c, 500, utils.RetrieveEntry, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
		)
	}

	password := c.Get(H.Conf.PASSWORD_HEADER_KEY)
	decryptedSecrets := []models.Secret{}

	for _, secret := range entry.Secrets {
		if decryptedString, err := utils.Decrypt(secret.String, password); err != nil {
			return utils.RespondWithError(c, 500, utils.RetrieveEntry, utils.ErrorDecrypt, err.Error())
		} else {
			secret.String = decryptedString
			decryptedSecrets = append(decryptedSecrets, secret)
		}
	}

	entry.Secrets = decryptedSecrets

	return c.Status(200).JSON(&entry)
}
