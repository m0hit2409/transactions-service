# Stage 1: build
FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# go-sqlite3 requires CGO; bookworm ships with gcc so no extra install needed.
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o api ./cmd/api

# Stage 2: minimal runtime
FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/api .

RUN mkdir -p data

EXPOSE 8080

CMD ["./api"]
