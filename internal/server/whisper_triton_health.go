package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const tritonReadyPath = "/v2/health/ready"

// whisperTritonHealth проверяет доступность контейнера Triton с model_repo Whisper
// (HTTP GET /v2/health/ready). Без WHISPER_TRITON_HTTP_URL проверка не выполняется.
func whisperTritonHealth() map[string]string {
	base := strings.TrimSpace(os.Getenv("WHISPER_TRITON_HTTP_URL"))
	out := make(map[string]string)
	if base == "" {
		out["whisper_triton_status"] = "not_configured"
		out["whisper_triton_running"] = "n/a"
		out["whisper_triton_healthy"] = "n/a"
		out["whisper_triton_message"] = "WHISPER_TRITON_HTTP_URL is not set"
		return out
	}
	base = strings.TrimRight(base, "/")
	url := base + tritonReadyPath

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		out["whisper_triton_status"] = "down"
		out["whisper_triton_running"] = "false"
		out["whisper_triton_healthy"] = "false"
		out["whisper_triton_message"] = fmt.Sprintf("unreachable: %v", err)
		return out
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		out["whisper_triton_status"] = "ready"
		out["whisper_triton_running"] = "true"
		out["whisper_triton_healthy"] = "true"
		out["whisper_triton_message"] = "Triton ready (/v2/health/ready)"
	default:
		out["whisper_triton_status"] = "not_ready"
		out["whisper_triton_running"] = "true"
		out["whisper_triton_healthy"] = "false"
		out["whisper_triton_message"] = fmt.Sprintf("Triton HTTP %d (not ready)", resp.StatusCode)
	}
	return out
}
