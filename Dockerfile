# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/seed ./cmd/seed

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget \
    && addgroup -S sasivision && adduser -S sasivision -G sasivision

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/seed .
COPY --from=builder /app/migrations ./migrations
COPY scripts/docker-entrypoint.sh ./docker-entrypoint.sh

RUN mkdir -p storage \
    && chmod +x docker-entrypoint.sh \
    && chown -R sasivision:sasivision /app

USER sasivision

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=40s --retries=3 \
  CMD wget -qO- http://127.0.0.1:${APP_PORT:-8080}/health || exit 1

ENV APP_ENV=production \
    APP_PORT=8080 \
    APP_NAME=SasiVision-API \
    RUN_MIGRATIONS=true \
    WAIT_FOR_DB=true \
    DB_WAIT_RETRIES=30 \
    MIGRATIONS_DIR=/app/migrations

ENTRYPOINT ["./docker-entrypoint.sh"]
