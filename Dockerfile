# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata \
    && addgroup -S sasivision && adduser -S sasivision -G sasivision

WORKDIR /app

COPY --from=builder /app/server .
RUN mkdir -p storage && chown -R sasivision:sasivision /app

USER sasivision

EXPOSE 8080

ENV APP_ENV=production \
    APP_PORT=8080 \
    APP_NAME=SasiVision-API

CMD ["./server"]
