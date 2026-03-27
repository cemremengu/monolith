FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache git nodejs npm
RUN GOBIN=/usr/local/bin go install github.com/go-task/task/v3/cmd/task@v3.49.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY web/package.json web/package-lock.json ./web/
RUN cd web && npm ci

COPY . .

RUN task build:linux

FROM debian:trixie-slim AS runtime

RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
RUN groupadd --gid 1001 monolith && useradd --uid 1001 --gid monolith --no-create-home --shell /bin/false monolith

WORKDIR /app

COPY --from=builder --chown=monolith:monolith /app/monolith /app/monolith

USER monolith

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/monolith"]
