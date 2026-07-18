package common

import "github.com/gofiber/fiber/v2"

// WriteJSON writes a successful JSON response with the given status code.
// Controllers should use this helper instead of calling c.JSON directly so
// response shaping stays consistent across the codebase.
func WriteJSON(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(data)
}
