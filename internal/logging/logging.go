// Package logging configures the process-wide slog logger from environment variables.
//
// LOG_LEVEL: debug | info | warn | error (default: info)
// LOG_FORMAT: text | json (default: text). Use json for log shipping (e.g. Filebeat → Elasticsearch / Kibana).
// LOG_VERBOSITY: standard | verbose (default: standard). Verbose adds structured request fields (query, User-Agent, sizes, …).
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

var verbosityVerbose bool

// InitFromEnv replaces [slog.Default] with a handler writing to w (typically os.Stderr).
func InitFromEnv(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	opts := &slog.HandlerOptions{Level: parseLevel(os.Getenv("LOG_LEVEL"))}

	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_VERBOSITY"))) {
	case "verbose":
		verbosityVerbose = true
	default:
		verbosityVerbose = false
	}

	var h slog.Handler
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT"))) {
	case "json":
		h = slog.NewJSONHandler(w, opts)
	default:
		h = slog.NewTextHandler(w, opts)
	}

	slog.SetDefault(slog.New(h))
}

// Verbose is true when LOG_VERBOSITY=verbose (HTTP middleware logs full field sets).
func Verbose() bool {
	return verbosityVerbose
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
