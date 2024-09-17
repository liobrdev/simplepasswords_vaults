package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

type UpdateSecretRequestBody struct {
	Label		 	string `json:"secret_label"`
	String	 	string `json:"secret_string"`
}

func (H Handler) UpdateSecret(c *fiber.Ctx) error {
	if clientOperation := c.Get("Client-Operation"); clientOperation == utils.MoveSecret {
		return c.Next()
	}

	body := UpdateSecretRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		return utils.RespondWithError(c, 400, utils.UpdateSecret, utils.ErrorParse, err.Error())
	}

	if body.Label == "" && body.String == "" {
		return utils.RespondWithError(
			c, 400, utils.UpdateSecret, utils.ErrorEmptyUpdateSecret, "Null or empty object or fields.",
		)
	}

	if len(body.Label) > 255 {
		return utils.RespondWithError(c, 400, utils.UpdateSecret, utils.ErrorSecretLabel, "Too long")
	}

	if len(body.String) > 1000 {
		return utils.RespondWithError(c, 400, utils.UpdateSecret, utils.ErrorSecretString, "Too long")
	}

	if body.String != "" {
		password := c.Get(H.Conf.PASSWORD_HEADER_KEY)

		var err error

		if body.String, err = utils.Encrypt(body.String, password); err != nil {
			return utils.RespondWithError(c, 500, utils.UpdateSecret, utils.ErrorEncrypt, err.Error())
		}
	}

	slug := c.Params("slug")

	if result := H.DB.Model(&models.Secret{}).
	Where("slug = ?", slug).Updates(models.Secret{Label: body.Label, String: body.String});
	result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.UpdateSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c, 404, utils.UpdateSecret, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c, 500, utils.UpdateSecret, "result.RowsAffected > 1", strconv.FormatInt(n, 10),
		)
	}

	return c.SendStatus(204)
}
