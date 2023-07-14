package user

import "github.com/gofiber/fiber/v2"

type UserController struct {
	storage *UserStorage
}

func NewUserController(storage *UserStorage) *UserController {
	return &UserController{
		storage: storage,
	}
}

type signUpRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Mobile    string `json:"mobile"`
}

func (u *UserController) register(c *fiber.Ctx) error {
	var req signUpRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	message, err := u.storage.signUp(req.FirstName, req.LastName, req.UserName, req.Mobile, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).SendString(message)
}
