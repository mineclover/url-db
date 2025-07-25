#!/bin/bash

# Docker build and push script for URL-DB

set -e

echo "=== Docker Build and Push Script ==="
echo

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
DOCKER_USERNAME="${DOCKER_USERNAME:-asfdassdssa}"
IMAGE_NAME="url-db"
FULL_IMAGE_NAME="${DOCKER_USERNAME}/${IMAGE_NAME}"

# Get version from git tag or use default
VERSION=$(git describe --tags --exact-match 2>/dev/null || echo "latest")

echo "Docker Hub Username: ${DOCKER_USERNAME}"
echo "Image Name: ${FULL_IMAGE_NAME}"
echo "Version: ${VERSION}"
echo

# Check if logged in to Docker Hub
if ! docker info | grep -q "Username: ${DOCKER_USERNAME}"; then
    echo -e "${YELLOW}Please log in to Docker Hub:${NC}"
    docker login
fi

# Build the image
echo -e "${GREEN}Building Docker image...${NC}"
docker build -t ${FULL_IMAGE_NAME}:${VERSION} .

# Tag as latest if not already
if [ "${VERSION}" != "latest" ]; then
    docker tag ${FULL_IMAGE_NAME}:${VERSION} ${FULL_IMAGE_NAME}:latest
fi

# Show image size
echo
echo -e "${GREEN}Built image details:${NC}"
docker images ${FULL_IMAGE_NAME}

# Test the image
echo
echo -e "${GREEN}Testing image...${NC}"

# Test stdio mode
echo "Testing stdio mode..."
echo '{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2025-06-18"},"id":1}' | \
    docker run -i --rm ${FULL_IMAGE_NAME}:latest -mcp-mode=stdio | \
    grep -q "result" && echo -e "${GREEN}✓ stdio mode test passed${NC}" || echo -e "${RED}✗ stdio mode test failed${NC}"

# Test SSE mode
echo "Testing SSE mode..."
docker run -d --name test-sse -p 8080:8080 ${FULL_IMAGE_NAME}:latest -mcp-mode=sse
sleep 3
if curl -s http://localhost:8080/health | grep -q "ok"; then
    echo -e "${GREEN}✓ SSE mode test passed${NC}"
else
    echo -e "${RED}✗ SSE mode test failed${NC}"
fi
docker stop test-sse && docker rm test-sse

# Ask for confirmation before pushing
echo
read -p "Do you want to push to Docker Hub? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}Pushing to Docker Hub...${NC}"
    
    # Push versioned tag
    docker push ${FULL_IMAGE_NAME}:${VERSION}
    
    # Push latest tag
    docker push ${FULL_IMAGE_NAME}:latest
    
    echo
    echo -e "${GREEN}✓ Successfully pushed to Docker Hub!${NC}"
    echo
    echo "Images available at:"
    echo "  - docker pull ${FULL_IMAGE_NAME}:latest"
    echo "  - docker pull ${FULL_IMAGE_NAME}:${VERSION}"
    echo
    echo "Run commands:"
    echo "  # stdio mode (for AI assistants)"
    echo "  docker run -it --rm -v \$(pwd)/data:/data ${FULL_IMAGE_NAME}:latest"
    echo
    echo "  # SSE mode (for HTTP clients)"
    echo "  docker run -d -p 8080:8080 -v \$(pwd)/data:/data --name url-db-sse ${FULL_IMAGE_NAME}:latest -mcp-mode=sse"
else
    echo "Push cancelled."
fi