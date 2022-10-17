package users

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type CreateUserRequestBody struct {
	Slug string `json:"user_slug"`
}

func (h handler) CreateUser(c *fiber.Ctx) error {
	body := CreateUserRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateUser,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if !utils.SlugRegexp.MatchString(body.Slug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.CreateUser,
			string(utils.ErrorUserSlug),
			body.Slug,
		)
	}

	if result := h.DB.Create(&models.User{Slug: body.Slug}); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.CreateUser,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.CreateUser,
			"result.RowsAffected != 1",
			strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
