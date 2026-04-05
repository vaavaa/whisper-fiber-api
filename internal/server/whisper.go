package server

import (
	"io"
	"mime/multipart"
	"time"
	"whisper-fiber-api/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WhisperHandler struct {
	db database.Service
}

// Transcribe enqueues a speech-to-text job in Redis.
// @Summary      Speech-to-text enqueue
// @Description  Accepts an audio file in the multipart form field `audio`, generates a `task_id` (UUID), and enqueues the job. Use WAV/MP3 or other formats supported by your downstream pipeline.
// @Tags         audio
// @Accept       multipart/form-data
// @Produce      json
// @Param        audio  formData  file  true  "Audio file (e.g. wav/mp3)"
// @Success      200  {object}  map[string]string  "JSON body includes task_id"
// @Failure      400  {string}  string  "Missing audio field or invalid multipart form"
// @Failure      500  {string}  string  "Failed to read the uploaded file or enqueue the job"
// @Router       /api/v1/transcribe [post]
func (h *WhisperHandler) Transcribe(c *fiber.Ctx) error {
	file, err := c.FormFile("audio")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("form field \"audio\" is required")
	}
	fileBytes, err := readFormFileBytes(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to read uploaded file")
	}

	taskID := uuid.NewString()
	ctx := c.UserContext()
	if err := h.db.EnqueueWhisperTask(ctx, taskID, fileBytes, time.Hour); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to enqueue task")
	}

	return c.JSON(fiber.Map{"task_id": taskID})
}

func readFormFileBytes(file *multipart.FileHeader) ([]byte, error) {
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
