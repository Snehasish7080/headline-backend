package thread

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ThreadController struct {
	storage *ThreadStorage
}

func NewThreadController(storage *ThreadStorage) *ThreadController {
	return &ThreadController{
		storage: storage,
	}
}

type threadRequest struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type threadResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (t *ThreadController) createThread(c *fiber.Ctx) error {
	var req threadRequest
	now := time.Now()
	uuid, uuidErr := uuid.NewRandom()

	opinionId := c.Params("id")

	if opinionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(threadResponse{
			Message: "invalid request",
			Success: false,
		})
	}

	if uuidErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(threadResponse{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(threadResponse{
			Message: "Invalid request body",
			Success: false,
		})
	}

	var threadFields map[string]interface{}
	data, _ := json.Marshal(req)
	json.Unmarshal(data, &threadFields)

	threadFields["created_at"] = now.Format(time.RFC3339)
	threadFields["updated_at"] = now.Format(time.RFC3339)
	threadFields["uuid"] = uuid.String()

	localData := c.Locals("userName")
	userName, cnvErr := localData.(string)

	if !cnvErr {
		return errors.New("not able to covert")
	}

	for k := range threadFields {
		if threadFields[k] == "" {
			delete(threadFields, k)
		}
	}

	message, err := t.storage.create(userName, opinionId, threadFields, c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(threadResponse{
			Message: "Creation Failed",
			Success: false,
		})

	}

	return c.Status(fiber.StatusOK).JSON(threadResponse{
		Message: message,
		Success: true,
	})

}

type allThread struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Created_at  string `json:"created_at"`
}

type getThreadResponse struct {
	Data    []*allThread `json:"data"`
	Message string       `json:"message"`
	Success bool         `json:"success"`
}

func (t *ThreadController) getAllThreads(c *fiber.Ctx) error {
	opinionId := c.Params("id")

	if opinionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(threadResponse{
			Message: "invalid request",
			Success: false,
		})
	}

	result, err := t.storage.get(opinionId, c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(getThreadResponse{
			Message: err.Error(),
			Success: false,
		})

	}

	jsonData, _ := json.Marshal(result)
	var structData []*allThread
	json.Unmarshal(jsonData, &structData)

	return c.Status(fiber.StatusOK).JSON(getThreadResponse{
		Data:    structData,
		Message: "found successfully",
		Success: true,
	})

}
