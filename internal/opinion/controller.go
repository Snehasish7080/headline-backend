package opinion

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OpinionController struct {
	storage *OpinionStorage
}

func NewOpinionController(storage *OpinionStorage) *OpinionController {
	return &OpinionController{
		storage: storage,
	}
}

type opinionRequest struct {
	Description string `json:"description"`
	Image       string `json:"image"`
}

type opinionResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (o *OpinionController) createOpinion(c *fiber.Ctx) error {
	var req opinionRequest
	now := time.Now()
	uuid, uuidErr := uuid.NewRandom()

	if uuidErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(opinionResponse{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(opinionResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	var opinionFields map[string]interface{}
	data, _ := json.Marshal(req)
	json.Unmarshal(data, &opinionFields)

	opinionFields["created_at"] = now.Format(time.RFC3339)
	opinionFields["updated_at"] = now.Format(time.RFC3339)
	opinionFields["uuid"] = uuid.String()

	localData := c.Locals("userName")
	userName, cnvErr := localData.(string)

	if !cnvErr {
		return errors.New("not able to covert")
	}

	for k := range opinionFields {
		if opinionFields[k] == "" {
			delete(opinionFields, k)
		}
	}

	message, err := o.storage.create(userName, opinionFields, c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(opinionResponse{
			Message: "Creation Failed",
			Success: false,
		})

	}

	return c.Status(fiber.StatusOK).JSON(opinionResponse{
		Message: message,
		Success: true,
	})
}
