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
		AllowHeaders:     "Accept,Authorization,Content-Type,Origin,X-Requested-With",
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))

	s.App.All("/echo", s.EchoHandler)
	s.App.All("/", s.EchoHandler)
	s.App.Get("/health", s.healthHandler)
	s.App.Get("/swagger/*", swagger.New(swagger.Config{InstanceName: "swagger"}))

	wh := &WhisperHandler{db: s.db}
	v1 := s.App.Group("/api/v1")
	v1.Post("/transcribe", wh.Transcribe)
}

// EchoHandler mirrors plain text for debugging.
// @Summary      Echo string
// @Description  Returns **`text/plain`**. If the request has a body, it is echoed as-is. If the body is empty, the query **`q`** is echoed when present; otherwise the response body is empty (HTTP 200, no error).
// @Tags         debug
// @Produce      plain
// @Param        q  query  string  false  "String to echo when the request body is empty (Swagger: use GET /echo)"
// @Success      200  {string}  string  "Request body bytes, or q, or empty"
// @Router       /echo [get]
func (s *FiberServer) EchoHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/plain; charset=utf-8")
	if body := c.Body(); len(body) > 0 {
		return c.Send(body)
	}
	q := c.Query("q")
	return c.Send([]byte(q))
}

// healthHandler checks Redis and related dependencies.
// @Summary      Health check
// @Description  Returns an aggregated dependency status (e.g. `redis_status`) and optional Triton/Whisper backend fields (`whisper_triton_status`, `whisper_triton_running`, `whisper_triton_healthy`, `whisper_triton_message`) when WHISPER_TRITON_HTTP_URL is set. Responds with HTTP 503 when Redis is unavailable.
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "All checks passed"
// @Failure      503  {object}  map[string]interface{}  "Service temporarily unavailable"
// @Router       /health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	stats := s.db.Health()
	for k, v := range whisperTritonHealth() {
		stats[k] = v
	}
	if stats["redis_status"] != "up" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(stats)
	}
	return c.JSON(stats)
}
