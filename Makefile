# Go 빌드 설정
BINARY_NAME=url-db
VERSION?=1.0.0
BUILD_DIR=build

# Go 컴파일러 설정
GO=go
GOFLAGS=-v

# 빌드 타깃 플랫폼
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# 빌드 플래그
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION}"

.PHONY: all build clean test deps run

# 기본 타겟
all: clean deps test build

# 의존성 설치
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# 빌드 (현재 플랫폼)
build:
	@echo "Building for current platform..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/server

# 멀티플랫폼 빌드
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1) \
		GOARCH=$$(echo $$platform | cut -d/ -f2) \
		output_name='$(BUILD_DIR)/$(BINARY_NAME)-'$$(echo $$platform | tr '/' '-'); \
		if [ $$GOOS = "windows" ]; then \
			output_name="$$output_name.exe"; \
		fi; \
		echo "Building $$output_name..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build $(GOFLAGS) $(LDFLAGS) -o $$output_name ./cmd/server; \
	done

# 실행
run: build
	./$(BINARY_NAME)

# 테스트
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# 테스트 (커버리지 포함)
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test ./... -v -cover -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

# 정적 분석
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: brew install golangci-lint"; \
	fi

# 포맷팅
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# 청소
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# 개발 모드 (hot reload)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Run: go install github.com/cosmtrek/air@latest"; \
	fi

# Swagger 문서 생성
swagger-gen:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o docs; \
	else \
		echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Swagger와 함께 개발 모드
dev-swagger: swagger-gen dev

# 도움말
help:
	@echo "Available targets:"
	@echo "  make deps          - Install dependencies"
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make run           - Build and run"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make dev           - Run in development mode with hot reload"
	@echo "  make swagger-gen   - Generate Swagger documentation"
	@echo "  make dev-swagger   - Generate Swagger docs and run dev mode"