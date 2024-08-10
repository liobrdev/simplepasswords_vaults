package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type CreateUserRequestBody struct {
	Slug string `json:"user_slug"`
}

func (H Handler) CreateUser(c *fiber.Ctx) error {
	body := CreateUserRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.CreateUser, utils.ErrorParse, err.Error())
	}

	if !utils.SlugRegexp.MatchString(body.Slug) {
		return utils.RespondWithError(c, 400, utils.CreateUser, utils.ErrorUserSlug, body.Slug)
	}

	if result := H.DB.Create(&models.User{ Slug: body.Slug }); result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: users.slug" {
			return utils.RespondWithError(
				c, 409, utils.CreateUser, utils.ErrorDuplicateUser, result.Error.Error(),
			)
		}

		return utils.RespondWithError(
			c, 500, utils.CreateUser, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n != 1 {
		return utils.RespondWithError(
			c, 500, utils.CreateUser, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(204)
}
