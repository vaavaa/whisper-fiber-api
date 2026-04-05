package server

import (
	"log/slog"
	"runtime/debug"

	"whisper-fiber-api/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type FiberServer struct {
	*fiber.App

	db  database.Service
	log *slog.Logger
}

func New() *FiberServer {
	log := slog.Default()
	app := fiber.New(fiber.Config{
		ServerHeader: "whisper-fiber-api",
		AppName:      "whisper-fiber-api",
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Error("panic recovered",
				"panic", e,
				"stack", string(debug.Stack()),
				"method", c.Method(),
				"path", c.Path(),
			)
		},
	}))
	app.Use(requestLogMiddleware(log))

	return &FiberServer{
		App: app,
		db:  database.New(),
		log: log,
	}
}
