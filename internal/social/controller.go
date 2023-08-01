package social

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type SocialController struct {
	storage *SocialStorage
}

func NewSocialController(storage *SocialStorage) *SocialController {
	return &SocialController{
		storage: storage,
	}
}

type followRequest struct {
	Account string `json:"account"`
}

type followResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (s *SocialController) followAccount(c *fiber.Ctx) error {
	var req followRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(followResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	localData := c.Locals("userName")
	userName, cnvErr := localData.(string)

	if !cnvErr {
		return errors.New("not able to covert")
	}

	if userName == req.Account {
		return c.Status(fiber.StatusBadRequest).JSON(followResponse{
			Message: "Invalid request",
			Success: false,
		})
	}

	message, err := s.storage.follow(userName, req.Account, c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(followResponse{
			Message: "Creation Failed",
			Success: false,
		})

	}

	return c.Status(fiber.StatusOK).JSON(followResponse{
		Message: message,
		Success: true,
	})
}
