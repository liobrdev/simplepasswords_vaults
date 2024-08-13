package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (H Handler) DeleteSecret(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(c, 400, utils.DeleteSecret, utils.ErrorSecretSlug, slug)
	}

	var secret models.Secret

	if result := H.DB.Delete(&secret, "slug = ?", slug); result.Error != nil {
		return utils.RespondWithError(
			c, 500, utils.DeleteSecret, utils.ErrorFailedDB, result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c, 404, utils.DeleteSecret, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c, 500, utils.DeleteSecret, fmt.Sprintf("result.RowsAffected (%d) > 1", n), "",
		)
	}

	return c.SendStatus(204)
}
