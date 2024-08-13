package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type UpdateSecretRequestBody struct {
	Label  string `json:"secret_label"`
	String string `json:"secret_string"`
}

func (H Handler) UpdateSecret(c *fiber.Ctx) error {
	body := UpdateSecretRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateSecret,
			string(utils.ErrorParse),
			err.Error(),
		)
	}

	if body.Label == "" && body.String == "" {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateSecret,
			string(utils.ErrorEmptyUpdateSecret),
			"Likely (null|empty) (object|fields).",
		)
	}

	if body.Label != "" {
		if len(body.Label) > 255 {
			return utils.RespondWithError(
				c,
				fiber.StatusBadRequest,
				utils.UpdateSecret,
				string(utils.ErrorSecretLabel),
				fmt.Sprintf("Too long (%d > 255)", len(body.Label)),
			)
		}
	}

	if body.String != "" && len(body.String) > 1000 {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.UpdateSecret,
			string(utils.ErrorSecretString),
			fmt.Sprintf("Too long (%d > 1000)", len(body.String)),
		)
	}

	slug := c.Params("slug")

	if result := H.DB.Model(&models.Secret{}).Where("slug = ?", slug).Updates(
		models.Secret{Label: body.Label, String: body.String},
	); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.UpdateSecret,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c,
			fiber.StatusNotFound,
			utils.UpdateSecret,
			string(utils.ErrorNoRowsAffected),
			"Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.UpdateSecret,
			"result.RowsAffected > 1",
			strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
