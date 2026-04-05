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

// Transcribe ставит задачу распознавания в очередь Redis.
// @Summary      Распознавание речи
// @Description  Принимает аудиофайл в поле формы `audio`, создаёт `task_id` (UUID) и ставит задачу в очередь. Формат: wav/mp3 и совместимые типы, поддерживаемые пайплайном.
// @Tags         audio
// @Accept       multipart/form-data
// @Produce      json
// @Param        audio  formData  file  true  "Аудиофайл (wav/mp3)"
// @Success      200  {object}  map[string]string  "task_id в теле JSON"
// @Failure      400  {string}  string  "Нет поля audio или неверная форма"
// @Failure      500  {string}  string  "Ошибка чтения файла или записи в очередь"
// @Router       /api/v1/transcribe [post]
func (h *WhisperHandler) Transcribe(c *fiber.Ctx) error {
	file, err := c.FormFile("audio")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("нужно поле формы audio")
	}
	fileBytes, err := readFormFileBytes(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("не удалось прочитать файл")
	}

	taskID := uuid.NewString()
	ctx := c.UserContext()
	if err := h.db.EnqueueWhisperTask(ctx, taskID, fileBytes, time.Hour); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("ошибка очереди")
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
