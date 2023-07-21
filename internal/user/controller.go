package user

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

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

type signUpResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (u *UserController) register(c *fiber.Ctx) error {
	var req signUpRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(signUpResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	token, err := u.storage.signUp(req.FirstName, req.LastName, req.UserName, req.Mobile, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(signUpResponse{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(signUpResponse{
		Token:   token,
		Success: true,
		Message: "Otp sent successfully",
	})
}

type verifyRequest struct {
	Otp string `json:"otp"`
}
type verifyResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (u *UserController) verifyOtp(c *fiber.Ctx) error {

	var req verifyRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(verifyResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	localData := c.Locals("userName")
	userName, cnvErr := localData.(string)

	if !cnvErr {
		return errors.New("not able to covert")
	}

	token, err := u.storage.verify(req.Otp, userName, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(verifyResponse{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(verifyResponse{
		Token:   token,
		Success: true,
		Message: "Verified",
	})

}

type loginRequest struct {
	Mobile string `json:"mobile"`
}
type loginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (u *UserController) loginUser(c *fiber.Ctx) error {
	var req loginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(loginResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	token, err := u.storage.login(req.Mobile, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(loginResponse{
			Message: err.Error(),
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(loginResponse{
		Token:   token,
		Success: true,
		Message: "Otp sent successfully",
	})
}

type userDetailResponse struct {
	Data    userDetail `json:"data"`
	Message string     `json:"message"`
	Success bool       `json:"success"`
}
type userDetail struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
}

func (u *UserController) getUserDetail(c *fiber.Ctx) error {
	localData := c.Locals("userName")
	userName, cnvErr := localData.(string)

	if !cnvErr {
		return errors.New("not able to covert")
	}
	user, err := u.storage.getUser(userName, c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(userDetailResponse{
			Message: err.Error(),
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(userDetailResponse{
		Data: userDetail{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			UserName:  user.UserName,
		},
		Message: "found successfully",
		Success: true,
	})
}

type updateUserDetailRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Image     string `json:"image"`
}
type updateUserDetailResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (u *UserController) updateUserDetail(c *fiber.Ctx) error {
	var req updateUserDetailRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(updateUserDetailResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	var updateFields map[string]interface{}
	data, _ := json.Marshal(req)
	json.Unmarshal(data, &updateFields)

	localData := c.Locals("userName")
	userName, cnvErr := localData.(string)

	if !cnvErr {
		return errors.New("not able to covert")
	}

	for k := range updateFields {
		if updateFields[k] == "" {
			delete(updateFields, k)
		}
	}

	message, err := u.storage.updateUser(userName, updateFields, c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(updateUserDetailResponse{
			Message: "Update failed",
			Success: false,
		})

	}

	return c.Status(fiber.StatusOK).JSON(updateUserDetailResponse{
		Message: message,
		Success: true,
	})
}
