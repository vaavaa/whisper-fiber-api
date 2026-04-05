// Whisper Fiber API — enqueue speech-to-text jobs for asynchronous processing.
// @title           Whisper Fiber API
// @version         1.0
// @description     HTTP API that queues Whisper-style audio transcription jobs via Redis. Versioned routes live under `/api/v1/...`.
// @BasePath        /
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"whisper-fiber-api/internal/logging"
	"whisper-fiber-api/internal/server"

	_ "github.com/joho/godotenv/autoload"
	_ "whisper-fiber-api/docs"
)

func gracefulShutdown(fiberServer *server.FiberServer, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := fiberServer.ShutdownWithContext(ctx); err != nil {
		slog.Error("server forced shutdown", "err", err)
	}

	slog.Info("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	logging.InitFromEnv(os.Stderr)

	server := server.New()

	server.RegisterFiberRoutes()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	go func() {
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		slog.Info("http server listening", "addr", fmt.Sprintf(":%d", port))
		if err := server.Listen(fmt.Sprintf(":%d", port)); err != nil {
			slog.Error("http server error", "err", err)
			panic(err)
		}
	}()

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("graceful shutdown complete")
}
