package server

import (
	"fmt"
	"log/slog"
	"time"

	"whisper-fiber-api/internal/logging"

	"github.com/gofiber/fiber/v2"
)

const httpReqMsg = "http request"

// requestLogMiddleware logs every HTTP request: Info for 2xx/3xx, Warn for 4xx, Error for 5xx or handler errors.
// With LOG_VERBOSITY=standard (default) only a single message line (method path status latency).
// With LOG_VERBOSITY=verbose, adds structured fields: route, ip, query, user_agent, content_type, bytes_in/out, err.
func requestLogMiddleware(log *slog.Logger) fiber.Handler {
	if log == nil {
		log = slog.Default()
	}
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		latency := time.Since(start)

		if logging.Verbose() {
			logHTTPVerbose(log, c, status, latency, err)
		} else {
			logHTTPStandard(log, c, status, latency, err)
		}
		return err
	}
}

func logHTTPStandard(log *slog.Logger, c *fiber.Ctx, status int, latency time.Duration, err error) {
	line := fmt.Sprintf("%s %s %d %dms", c.Method(), c.Path(), status, latency.Milliseconds())
	if err != nil {
		line = fmt.Sprintf("%s: %v", line, err)
	}
	switch {
	case err != nil && status >= 400 && status < 500:
		log.Warn(line)
	case err != nil:
		log.Error(line)
	case status >= 500:
		log.Error(line)
	case status >= 400:
		log.Warn(line)
	default:
		log.Info(line)
	}
}

func logHTTPVerbose(log *slog.Logger, c *fiber.Ctx, status int, latency time.Duration, err error) {
	attrs := []any{
		"method", c.Method(),
		"path", c.Path(),
		"route", c.Route().Path,
		"status", status,
		"latency_ms", latency.Milliseconds(),
		"ip", c.IP(),
		"query", c.Queries(),
		"user_agent", c.Get(fiber.HeaderUserAgent),
		"content_type", c.Get(fiber.HeaderContentType),
		"bytes_in", len(c.Body()),
		"bytes_out", len(c.Response().Body()),
	}
	if err != nil {
		attrs = append(attrs, "err", err)
	}
	switch {
	case err != nil && status >= 400 && status < 500:
		log.Warn(httpReqMsg, attrs...)
	case err != nil:
		log.Error(httpReqMsg, attrs...)
	case status >= 500:
		log.Error(httpReqMsg, attrs...)
	case status >= 400:
		log.Warn(httpReqMsg, attrs...)
	default:
		log.Info(httpReqMsg, attrs...)
	}
}
