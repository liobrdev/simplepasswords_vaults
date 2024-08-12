package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) RetrieveUser(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.RetrieveUser, utils.ErrorUserSlug, slug)
	}

	var user models.User

	if result := H.DB.Preload("Vaults").First(&user, "slug = ?", slug); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.RespondWithError(c, 404, utils.RetrieveUser, utils.ErrorNotFound, slug)
		}

		return utils.RespondWithError(
			c, 500, utils.RetrieveUser, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c, 500, utils.RetrieveUser, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
		)
	}

	if len(user.Vaults) == 0 {
		user.Vaults = []models.Vault{}
	}

	return c.Status(200).JSON(&user)
}
