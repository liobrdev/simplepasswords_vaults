package utils

import "github.com/gofiber/fiber/v2"

func RespondWithError(
	c *fiber.Ctx,
	statusCode int,
	operation string,
	message string,
	detail string,
) error {
	return c.Status(statusCode).JSON(ErrorResponseBody{
		ClientOperation: operation,
		Message:         message,
		ContextString:   c.String(),
		RequestBody:     string(c.Body()),
		Detail:          detail,
	})
}
