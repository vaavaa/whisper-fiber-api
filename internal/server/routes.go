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

// EchoHandler mirrors the request for debugging.
// @Summary      Echo request
// @Description  **GET** returns query parameters as a JSON object, or `{}` if empty. **Other methods** (POST, etc.) return the request body; when the body is non-empty, the original `Content-Type` is echoed back.
// @Tags         debug
// @Produce      json
// @Success      200  {object}  map[string]string  "GET: query keys and values"
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

// healthHandler checks Redis and related dependencies.
// @Summary      Health check
// @Description  Returns an aggregated dependency status (e.g. `redis_status`). Responds with HTTP 503 when Redis is unavailable.
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "All checks passed"
// @Failure      503  {object}  map[string]interface{}  "Service temporarily unavailable"
// @Router       /health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	stats := s.db.Health()
	if stats["redis_status"] != "up" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(stats)
	}
	return c.JSON(stats)
}
