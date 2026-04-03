package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	wh := &WhisperHandler{db: s.db}
	s.App.Post("/transcribe", wh.Transcribe)
}

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

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	stats := s.db.Health()
	if stats["redis_status"] != "up" {
		return c.Status(fiber.StatusServiceUnavailable).JSON(stats)
	}
	return c.JSON(stats)
}
