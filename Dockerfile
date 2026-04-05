# syntax=docker/dockerfile:1

FROM golang:bookworm AS builder

WORKDIR /src

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /api /api

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/api"]
