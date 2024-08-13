package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type UpdateEntryRequestBody struct {
	Title string `json:"entry_title"`
}

func (H Handler) UpdateEntry(c *fiber.Ctx) error {
	body := UpdateEntryRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.UpdateEntry, utils.ErrorParse, err.Error())
	}

	if body.Title == "" {
		return utils.RespondWithError(c, 400, utils.UpdateEntry, utils.ErrorEntryTitle, "")
	}

	if len(body.Title) > 255 {
		return utils.RespondWithError(c, 400, utils.UpdateEntry, utils.ErrorEntryTitle, "Too long")
	}

	slug := c.Params("slug")

	if result := H.DB.Model(&models.Entry{}).Where("slug = ?", slug).Update("title", body.Title);
	result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.UpdateEntry, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c, 404, utils.UpdateEntry, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c, 500, utils.UpdateEntry, "result.RowsAffected > 1", strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(204)
}
