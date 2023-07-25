package thread

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zone/headline/internal/middleware"
)

func AddThreadRoutes(app *fiber.App, middleware *middleware.AuthMiddleware, controller *ThreadController) {
	opinion := app.Group("/thread", middleware.VerifyUser)

	// add routes here
	opinion.Post("/:id", controller.createThread)

}
