package extensions

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// RegisterOpenAPI mounts the OpenAPI UI (swaggo) at /openapi/*.
// Regenerate src/docs with: make openapi
func RegisterOpenAPI(app *fiber.App) {
	app.Get("/openapi/*", fiberSwagger.WrapHandler)
}
