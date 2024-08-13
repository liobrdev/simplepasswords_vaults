package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) DeleteEntry(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.DeleteEntry, utils.ErrorEntrySlug, slug)
	}

	var entry models.Entry
	var result *gorm.DB

	if err := H.DB.Transaction(func(tx *gorm.DB) error {
		if result = tx.Delete(&entry, "slug = ?", slug); result.Error != nil {
			return result.Error
		} else if n := result.RowsAffected; n == 0 {
			return errors.New(utils.ErrorNoRowsAffected)
		} else if n > 1 {
			return fmt.Errorf("result.RowsAffected (%d) > 1", n)
		}

		if result = tx.Delete(&models.Secret{}, "entry_slug = ?", slug); result.Error != nil {
			return result.Error
		}

		return nil
	}); err != nil {
		if errText := err.Error(); errText == utils.ErrorNoRowsAffected {
			return utils.RespondWithError(
				c, 404, utils.DeleteEntry, errText, "Likely that slug was not found.",
			)
		} else {
			return utils.RespondWithError(c, 500, utils.DeleteEntry, utils.ErrorFailedDB, err.Error())
		}
	}

	return c.SendStatus(204)
}
