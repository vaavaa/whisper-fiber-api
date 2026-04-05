# whisper-fiber-api

**Turn raw audio into structured work—not into blocked HTTP threads.** This service is a fast, production-minded entry point for Whisper-style speech recognition: upload a clip, get a task id back immediately, and let Redis-backed streaming queue the heavy lifting for your workers and NVIDIA Triton.

Under the hood you get a lean **[Fiber](https://gofiber.io/)** API in Go, **[Redis](https://redis.io/)** for short-lived audio storage and a **`whisper_tasks`** stream for downstream consumers, plus deployment assets and a **gRPC Triton client** skeleton for Whisper ensemble inference (`internal/tritonwhisper`).

---

## What’s in this repo

- **HTTP API** — `POST /api/v1/transcribe` accepts multipart form field `audio`, stores bytes in Redis, publishes `task_id` to a Redis Stream, returns `{ "task_id": "..." }`.
- **Health** — `GET /health` reports Redis connectivity (503 if Redis is down).
- **Playground** — `GET|POST /` echoes body or query (useful for smoke tests).
- **Inference stack (templates)** — `deployments/` includes a **NVIDIA Triton** `docker-compose` and a **model repository** layout (ensemble + Python preprocess); weights are expected locally per `.gitignore`, not committed.
- **Go Triton client** — gRPC helpers for ensemble / decoder calls via [`go-triton-client`](https://github.com/Trendyol/go-triton-client).

> **Note:** This repository implements the **API producer** side of the pipeline (enqueue). A separate consumer should read `whisper_tasks`, fetch audio by key, call Triton (or another engine), and persist results— wired to your product.

---

## Requirements

- **Go** 1.26+ (see `go.mod`)
- **Redis** 7.x (or compatible) for production; local dev can use Docker Compose in this repo

---

## Quick start

1. **Start Redis**

   ```bash
   make docker-run
   ```

   Stops with `make docker-down`.

2. **Configure environment** (see variables below). A minimal local set typically includes Redis address, port, DB index, and `PORT` for the API.

3. **Run the API**

   ```bash
   make run
   ```

   Or build and run the binary:

   ```bash
   make build
   ./main
   ```

4. **(Optional) Triton** — from `deployments/`, follow comments in `docker-compose.triton.yml` to run the model repository on GPU or CPU. You still need model artifacts (e.g. ONNX) as described there.

---

## API

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/transcribe` | Form field **`audio`**: file upload. Returns JSON with **`task_id`**. |
| `GET` | `/health` | JSON health; **503** if Redis is unavailable. |
| `*`| `/` | Echo handler for quick checks. |

---

## OpenAPI / Swagger

Interactive docs are served by the app at **`/swagger/index.html`** (or open **`/swagger/`**).

Spec files in the repo are generated with [swag](https://github.com/swaggo/swag). After you change Swag comments in Go sources, regenerate them from the project root:

```bash
swag init -g cmd/api/main.go
```

That refreshes `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`. You may see a harmless warning if there are no `.go` files in the module root; generation still succeeds.

---

## Environment variables

| Variable | Role |
|----------|------|
| `PORT` | HTTP listen port for Fiber. |
| `BLUEPRINT_DB_ADDRESS` | Redis host. |
| `BLUEPRINT_DB_PORT` | Redis port (Compose maps host port → container 6379). |
| `BLUEPRINT_DB_PASSWORD` | Redis password (empty if none). |
| `BLUEPRINT_DB_DATABASE` | Redis **logical DB index** (string parsed as int). |

[`godotenv`](https://github.com/joho/godotenv) auto-loads a `.env` file if present.

---

## Makefile

| Target | Description |
|--------|-------------|
| `make all` | `build` + `test` |
| `make build` | Build `main` from `cmd/api/main.go` |
| `make run` | `go run cmd/api/main.go` |
| `make test` | `go test ./... -v` |
| `make itest` | Integration tests (`internal/database`) |
| `make docker-run` / `make docker-down` | Redis stack via Docker Compose |
| `make watch` | Live reload with [Air](https://github.com/air-verse/air) (prompts to install if missing) |
| `make clean` | Remove built `main` |
