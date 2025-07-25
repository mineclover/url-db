# Go 빌드 설정
BINARY_NAME=url-db
VERSION?=1.0.0
BUILD_DIR=bin

# Go 컴파일러 설정
GO=go
GOFLAGS=-v

# 빌드 타깃 플랫폼
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# 빌드 플래그
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION}"

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m

.PHONY: all build clean deps run build-all lint fmt dev swagger-gen dev-swagger help test test-coverage coverage-analysis docker-build docker-run docker-stop docker-logs docker-push

# 기본 타겟
all: clean deps build

# 의존성 설치
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

# 빌드 (현재 플랫폼)
build:
	@echo "$(BLUE)Building URL-DB Server...$(NC)"
	@echo "$(BLUE)Building for current platform...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go
	@if [ $$? -ne 0 ]; then \
		echo "$(RED)✗ Build failed!$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✓ Build completed successfully!$(NC)"
	@echo "$(GREEN)✓ Executable created: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# 멀티플랫폼 빌드
build-all:
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1) \
		GOARCH=$$(echo $$platform | cut -d/ -f2) \
		output_name='$(BUILD_DIR)/$(BINARY_NAME)-'$$(echo $$platform | tr '/' '-'); \
		if [ $$GOOS = "windows" ]; then \
			output_name="$$output_name.exe"; \
		fi; \
		echo "$(BLUE)Building $$output_name...$(NC)"; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build $(GOFLAGS) $(LDFLAGS) -o $$output_name cmd/server/main.go; \
		if [ $$? -ne 0 ]; then \
			echo "$(RED)✗ Build failed for $$platform!$(NC)"; \
			exit 1; \
		fi; \
	done
	@echo "$(GREEN)✓ Multi-platform build completed successfully!$(NC)"

# 실행
run: build
	@echo "$(BLUE)Running the server...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME)

# 정적 분석
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
		if [ $$? -eq 0 ]; then \
			echo "$(GREEN)✓ Linting passed$(NC)"; \
		else \
			echo "$(RED)✗ Linting failed$(NC)"; \
			exit 1; \
		fi; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not installed. Run: brew install golangci-lint$(NC)"; \
	fi

# 포맷팅
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

# 청소
clean:
	@echo "$(BLUE)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Clean completed$(NC)"

# 개발 모드 (hot reload)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(YELLOW)⚠ air not installed. Run: go install github.com/cosmtrek/air@latest$(NC)"; \
	fi

# Swagger 문서 생성
swagger-gen:
	@echo "$(BLUE)Generating Swagger documentation...$(NC)"
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o docs; \
		echo "$(GREEN)✓ Swagger documentation generated$(NC)"; \
	else \
		echo "$(YELLOW)⚠ swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest$(NC)"; \
	fi

# Swagger와 함께 개발 모드
dev-swagger: swagger-gen dev

# 테스트 실행
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@./scripts/test_runner.sh
	@echo "$(GREEN)✓ Tests completed$(NC)"

# 테스트 + 커버리지
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@./scripts/test_runner.sh -m coverage
	@echo "$(GREEN)✓ Tests with coverage completed$(NC)"

# 커버리지 분석만
coverage-analysis:
	@echo "$(BLUE)Running coverage analysis...$(NC)"
	@./scripts/coverage_analysis.sh
	@echo "$(GREEN)✓ Coverage analysis completed$(NC)"

# 도움말
help:
	@echo "$(BLUE)Available targets:$(NC)"
	@echo "  make deps          - Install dependencies"
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make run           - Build and run"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make dev           - Run in development mode with hot reload"
	@echo "  make swagger-gen   - Generate Swagger documentation"
	@echo "  make dev-swagger   - Generate Swagger docs and run dev mode"
	@echo ""
	@echo "$(BLUE)Testing commands:$(NC)"
	@echo "  make test              - Run all tests"
	@echo "  make test-coverage     - Run tests with detailed coverage analysis"
	@echo "  make coverage-analysis - Run coverage analysis only"
	@echo "  ./scripts/test_runner.sh -h  - Show test runner options"
	@echo "  ./scripts/coverage_analysis.sh - Detailed coverage analysis script"
	@echo ""
	@echo "$(BLUE)Docker commands:$(NC)"
	@echo "  make docker-build      - Build Docker image"
	@echo "  make docker-run        - Run container in MCP stdio mode"
	@echo "  make docker-compose-up - Start all services with Docker Compose"
	@echo "  make docker-compose-down - Stop all services"
	@echo "  make docker-logs       - Show Docker logs"
	@echo "  make docker-push       - Push image to registry"
	@echo "  make docker-clean      - Clean Docker resources"
	@echo ""
	@echo "$(BLUE)To run the server:$(NC)"
	@echo "  ./$(BUILD_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "$(BLUE)Default configuration (managed in internal/constants/):$(NC)"
	@echo "  Port: 8080 (constants.DefaultPort)"
	@echo "  Database: file:./url-db.sqlite (constants.DefaultDBPath)"
	@echo "  Tool Name: $(BINARY_NAME) (constants.DefaultServerName)"
	@echo "  MCP Server: url-db-mcp-server (constants.MCPServerName)"

# Docker 관련 설정
DOCKER_IMAGE=url-db
DOCKER_TAG?=latest
DOCKER_REGISTRY?=

# Docker 빌드
docker-build:
	@echo "$(BLUE)Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@if [ $$? -eq 0 ]; then \
		echo "$(GREEN)✓ Docker image built successfully!$(NC)"; \
		echo "$(GREEN)✓ Image: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"; \
	else \
		echo "$(RED)✗ Docker build failed!$(NC)"; \
		exit 1; \
	fi

# Docker 실행 (MCP stdio mode)
docker-run:
	@echo "$(BLUE)Running Docker container in MCP stdio mode...$(NC)"
	@echo "$(YELLOW)Use Ctrl+C to stop the container$(NC)"
	docker run -it --rm \
		--name url-db-mcp \
		-v url-db-data:/data \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker Compose 실행 (모든 서비스)
docker-compose-up:
	@echo "$(BLUE)Starting all services with Docker Compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✓ Services started!$(NC)"
	@echo "$(GREEN)✓ HTTP API: http://localhost:8080$(NC)"
	@echo "$(GREEN)✓ MCP SSE: http://localhost:8081$(NC)"
	@echo "$(GREEN)✓ MCP HTTP: http://localhost:8082$(NC)"
	@echo "$(YELLOW)Run 'make docker-logs' to see logs$(NC)"

# Docker Compose 중지
docker-compose-down:
	@echo "$(BLUE)Stopping all services...$(NC)"
	docker-compose down
	@echo "$(GREEN)✓ Services stopped$(NC)"

# Docker 로그 보기
docker-logs:
	@echo "$(BLUE)Showing Docker logs...$(NC)"
	docker-compose logs -f

# Docker 이미지 푸시
docker-push:
	@if [ -z "$(DOCKER_REGISTRY)" ]; then \
		echo "$(RED)✗ DOCKER_REGISTRY not set!$(NC)"; \
		echo "$(YELLOW)Usage: make docker-push DOCKER_REGISTRY=your-registry.com$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Tagging image for registry...$(NC)"
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(BLUE)Pushing to registry...$(NC)"
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)✓ Image pushed to $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

# Docker 정리
docker-clean:
	@echo "$(BLUE)Cleaning Docker resources...$(NC)"
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	@echo "$(GREEN)✓ Docker resources cleaned$(NC)"