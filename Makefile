# Makefile — сборка Go API, Docker (Redis + API), отдельно Triton с Whisper model_repo.
#
# Кратко:
#   make build|run|test   — локальная сборка и тесты
#   make docker-run       — только Redis (удобно для make run снаружи)
#   make docker-up        — API + Redis (корневой docker-compose.yml)
#   make docker-down      — остановить корневой compose
#   make triton-up-gpu    — Triton с GPU (нужны драйвер и NVIDIA Container Toolkit)
#   make triton-up-cpu    — Triton без GPU (профиль cpu; см. предупреждение в deployments/docker-compose.triton.yml)
#   make triton-down      — остановить compose Triton в deployments/

.PHONY: all build run test clean watch docker-run docker-up docker-down itest \
	triton-up-gpu triton-up-cpu triton-down

# Сборка бинарника API
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Запуск API на хосте (Redis поднимите отдельно: make docker-run)
run:
	@go run cmd/api/main.go

# Только Redis из корневого compose — для локального `make run` против Docker Redis
docker-run:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose up -d redis_bp; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up -d redis_bp; \
	else \
		echo "Docker Compose не найден. Установите плагин Compose V2 (docker compose)."; \
		exit 1; \
	fi

# API + Redis (образ собирается из Dockerfile в корне)
docker-up:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose up --build; \
	else \
		echo "Docker Compose не найден. Установите плагин Compose V2 (docker compose)."; \
		exit 1; \
	fi

# Остановить сервисы корневого docker-compose.yml
docker-down:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose down; \
	else \
		echo "Docker Compose не найден. Установите плагин Compose V2 (docker compose)."; \
		exit 1; \
	fi

# Triton + Whisper: GPU, порты по умолчанию 8000/8001/8002 (см. TRITON_*_PORT в .env)
# Запуск из deployments/ — путь ./model_repo в compose остаётся корректным.
triton-up-gpu:
	@cd deployments && \
	if docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.triton.yml up -d --pull always; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f docker-compose.triton.yml up -d --pull always; \
	else \
		echo "Docker Compose не найден. Установите плагин Compose V2 (docker compose)."; \
		exit 1; \
	fi

# Triton без GPU (--profile cpu), порты 8100/8101/8102 по умолчанию.
# В model_repo для CPU может понадобиться KIND_CPU в config.pbtxt — см. комментарии в deployments/docker-compose.triton.yml
triton-up-cpu:
	@cd deployments && \
	if docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.triton.yml --profile cpu up -d --pull always; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f docker-compose.triton.yml --profile cpu up -d --pull always; \
	else \
		echo "Docker Compose не найден. Установите плагин Compose V2 (docker compose)."; \
		exit 1; \
	fi

# Остановить сервисы из deployments/docker-compose.triton.yml
triton-down:
	@cd deployments && \
	if docker compose version >/dev/null 2>&1; then \
		docker compose -f docker-compose.triton.yml down; \
	elif command -v docker-compose >/dev/null 2>&1; then \
		docker-compose -f docker-compose.triton.yml down; \
	else \
		echo "Docker Compose не найден. Установите плагин Compose V2 (docker compose)."; \
		exit 1; \
	fi

# Тесты
test:
	@echo "Testing..."
	@go test ./... -v

# Интеграционные тесты (Redis через testcontainers)
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

clean:
	@echo "Cleaning..."
	@rm -f main

# Live reload (требуется air)
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching..."; \
	else \
		read -p "Установить air? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching..."; \
		else \
			echo "Выход без air."; \
			exit 1; \
		fi; \
	fi
