package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zone/headline/internal/middleware"
)

func AddUserRoutes(app *fiber.App, middleware *middleware.AuthMiddleware, controller *UserController) {
	auth := app.Group("/auth")

	// add routes here
	auth.Post("/sign-up", controller.register)
	auth.Post("/login", controller.loginUser)

	verify := auth.Group("/verify", middleware.VerifyOtpToken)
	verify.Get("/", controller.verifyOtp)

}
