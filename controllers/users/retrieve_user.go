package users

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (h handler) RetrieveUser(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.RetrieveUser,
			string(utils.ErrorUserSlug),
			slug,
		)
	}

	var user models.User

	if result := h.DB.Preload("Vaults").First(&user, "slug = ?", slug); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return utils.RespondWithError(
				c,
				fiber.StatusNotFound,
				utils.RetrieveUser,
				string(utils.ErrorNotFound),
				slug,
			)
		}

		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.RetrieveUser,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.RetrieveUser,
			"result.RowsAffected != 1",
			strconv.FormatInt(n, 10),
		)
	}

	return c.Status(fiber.StatusOK).JSON(&user)
}
