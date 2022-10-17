package entries

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type UpdateEntryRequestBody struct {
	Title string `json:"entry_title"`
}

func (h handler) UpdateEntry(c *fiber.Ctx) error {
	body := UpdateEntryRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateEntry,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if body.Title == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateEntry,
			string(utils.ErrorEntryTitle),
			"",
		)
	}

	if len(body.Title) > 255 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateEntry,
			string(utils.ErrorEntryTitle),
			fmt.Sprintf("Too long (%d > 255)", len(body.Title)),
		)
	}

	slug := c.Params("slug")

	if result := h.DB.Model(&models.Entry{}).Where("slug = ?", slug).Update(
		"title",
		body.Title,
	); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.UpdateEntry,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c,
			fiber.StatusNotFound,
			utils.UpdateEntry,
			string(utils.ErrorNoRowsAffected),
			"Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.UpdateEntry,
			"result.RowsAffected > 1",
			strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
