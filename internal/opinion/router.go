package opinion

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zone/headline/internal/middleware"
)

func AddOpinionRoutes(app *fiber.App, middleware *middleware.AuthMiddleware, controller *OpinionController) {
	opinion := app.Group("/opinion", middleware.VerifyUser)

	// add routes here
	opinion.Post("/", controller.createOpinion)

}
