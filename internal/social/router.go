package social

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zone/headline/internal/middleware"
)

func AddSocialRoutes(app *fiber.App, middleware *middleware.AuthMiddleware, controller *SocialController) {
	opinion := app.Group("/follow", middleware.VerifyUser)

	// add routes here
	opinion.Post("/", controller.followAccount)

}
