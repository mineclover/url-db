#!/bin/bash

# Simple SSE mode setup script

set -e

echo "=== URL-DB SSE Mode Setup ==="
echo

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check Docker
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is required but not installed"
    echo "Install Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# Create data directory
mkdir -p data
echo -e "${GREEN}✓ Data directory created${NC}"

echo
echo "=== Setup Complete ==="
echo
echo "Start SSE server:"
echo "  docker run -d -p 8080:8080 -v \$(pwd)/data:/data --name url-db-sse asfdassdssa/url-db:latest -mcp-mode=sse"
echo
echo "Or with Docker Compose:"
echo "  docker-compose -f docker-compose-sse.yml up -d"
echo
echo "Test connection:"
echo "  curl http://localhost:8080/health"
echo
echo "Stop server:"
echo "  docker stop url-db-sse && docker rm url-db-sse"
echo "  # or: docker-compose -f docker-compose-sse.yml down"
echo