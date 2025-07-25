# Stage 1: Build web frontend
FROM node:22-alpine AS web-builder

WORKDIR /app/web

# Copy package files first for better caching
COPY web/package*.json ./ 
RUN npm ci --only=production

# Copy web source and build
COPY web/ ./
RUN npm run build

# Stage 2: Build Go application
FROM golang:1.24-alpine AS go-builder

# Install git for go modules that might need it
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built web assets from previous stage
COPY --from=web-builder /app/web/dist ./web/dist

# Build the Go application
ARG VERSION=v0.0.0
ARG COMMIT=unknown
ARG DATE_BUILT
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build \
  -ldflags "-w -s -X monolith.Version=${VERSION} -X monolith.Commit=${COMMIT} -X monolith.DateBuilt=${DATE_BUILT}" \
  -o monolith \
  ./cmd/monolith

# Stage 3: Final runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S monolith && \
  adduser -u 1001 -S monolith -G monolith

WORKDIR /app

# Copy the binary from builder stage
COPY --from=go-builder /app/monolith .

# Copy any static files if needed (adjust path as necessary)
# COPY --from=go-builder /app/static ./static

# Change ownership to non-root user
RUN chown -R monolith:monolith /app
USER monolith

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run
ENTRYPOINT ["./monolith"]
