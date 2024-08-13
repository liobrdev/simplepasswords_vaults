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
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.RetrieveEntry,
			string(utils.ErrorEntrySlug),
			slug,
		)
	}

	var entry models.Entry

	if result := H.DB.Preload("Secrets").First(&entry, "slug = ?", slug); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.RespondWithError(
				c,
				fiber.StatusNotFound,
				utils.RetrieveEntry,
				string(utils.ErrorNotFound),
				slug,
			)
		}

		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.RetrieveEntry,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.RetrieveEntry,
			"result.RowsAffected != 1",
			strconv.FormatInt(n, 10),
		)
	}

	return c.Status(fiber.StatusOK).JSON(&entry)
}
