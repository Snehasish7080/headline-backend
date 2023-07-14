package user

import "github.com/gofiber/fiber/v2"

func AddUserRoutes(app *fiber.App, controller *UserController) {
	auth := app.Group("/auth")

	// add routes here
	auth.Post("/sign-up", controller.register)

}
