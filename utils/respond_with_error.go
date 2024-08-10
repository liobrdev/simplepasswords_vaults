package utils

import (
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RespondWithError(c *fiber.Ctx, statusCode int, operation, message, detail string) error {
	_, file, line, _ := runtime.Caller(1)

	return c.Status(statusCode).JSON(ErrorResponseBody{
		Caller:					 file + ":" + strconv.FormatInt(int64(line), 10),
		ClientOperation: operation,
		Message:         message,
		ContextString:   c.String(),
		RequestBody:     string(c.Body()),
		Detail:          detail,
	})
}
