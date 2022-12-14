package secrets

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func (h handler) DeleteSecret(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if !utils.SlugRegexp.MatchString(slug) {
		return utils.RespondWithError(
			c,
			fiber.StatusBadRequest,
			utils.DeleteSecret,
			string(utils.ErrorSecretSlug),
			slug,
		)
	}

	var secret models.Secret

	if result := h.DB.Delete(&secret, "slug = ?", slug); result.Error != nil {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.DeleteSecret,
			string(utils.ErrorFailedDB),
			result.Error.Error(),
		)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(
			c,
			fiber.StatusNotFound,
			utils.DeleteSecret,
			string(utils.ErrorNoRowsAffected),
			"Likely that slug was not found.",
		)
	} else if n > 1 {
		return utils.RespondWithError(
			c,
			fiber.StatusConflict,
			utils.DeleteSecret,
			fmt.Sprintf("result.RowsAffected (%d) > 1", n),
			"",
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
