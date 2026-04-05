package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func (s *FiberServer) RegisterFiberRoutes() {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.All("/", s.EchoHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/swagger/*", swagger.HandlerDefault)

	wh := &WhisperHandler{db: s.db}
	v1 := s.App.Group("/api/v1")
	v1.Post("/transcribe", wh.Transcribe)
}

// EchoHandler отражает запрос для отладки.
// @Summary      Эхо запроса
// @Description  **GET** — возвращает query-параметры JSON-объектом или `{}`. **Остальные методы** (POST и др.) — возвращают тело запроса; при непустом теле повторяется исходный `Content-Type`.
// @Tags         debug
// @Produce      json
// @Success      200  {object}  map[string]string  "GET: ключи из query"
// @Router       / [get]
func (s *FiberServer) EchoHandler(c *fiber.Ctx) error {
	body := c.Body()
	if len(body) > 0 {
		if ct := c.Get("Content-Type"); ct != "" {
			c.Set("Content-Type", ct)
		}
		return c.Send(body)
	}
	if q := c.Queries(); len(q) > 0 {
		return c.JSON(q)
	}
	return c.JSON(fiber.Map{})
}

// healthHandler проверяет доступность Redis и связанных зависимостей.
// @Summary      Проверка готовности
// @Description  Возвращает агрегированный статус сервисов (например, `redis_status`). При недоступности Redis — HTTP 503.
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Все проверки прошли"
// @Failure      503  {object}  map[string]interface{}  "Сервис временно недоступен"
// @Router       /health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	stats := s.db.Health()
	if stats["redis_status"] != "up" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(stats)
	}
	return c.JSON(stats)
}
