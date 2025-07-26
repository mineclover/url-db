# Multi-stage build for optimal image size
FROM golang:1.24.5-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w -X main.Version=1.0.0" -o url-db cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o mcp-bridge cmd/bridge/main.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite-libs

# Create non-root user
RUN addgroup -g 1000 -S urldb && \
    adduser -u 1000 -S urldb -G urldb

# Set working directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/url-db .
COPY --from=builder /build/mcp-bridge .

# Copy schema file
COPY --from=builder /build/schema.sql .

# Create directory for database with proper permissions
RUN mkdir -p /data && chown -R urldb:urldb /data /app

# Switch to non-root user
USER urldb

# Volume for persistent data
VOLUME ["/data"]

# Expose ports for different modes
# 8080 for HTTP/SSE mode
EXPOSE 8080

# Default environment variables
ENV DATABASE_URL="file:/data/url-db.sqlite"
ENV TOOL_NAME="url-db"

# Default command for MCP stdio mode
ENTRYPOINT ["./url-db"]
CMD ["-mcp-mode=stdio", "-db-path=/data/url-db.sqlite"]