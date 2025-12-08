# LexVeritas Backend Makefile
# ============================================================================

# 变量
APP_NAME := lex-veritas-backend
BUILD_DIR := bin
MAIN_FILE := cmd/server/main.go
MIGRATE_FILE := cmd/migrate/main.go
CONFIG_FILE := config.yaml

# Go 命令
GO := go
GOFLAGS := -v

# 版本信息
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

.PHONY: all build run dev test clean migrate swagger help

# ============================================================================
# 默认目标
# ============================================================================

all: help

# ============================================================================
# 开发命令
# ============================================================================

## dev: 使用 Air 热重载运行
dev:
	@echo "🚀 Starting development server with hot reload..."
	@air

## run: 直接运行服务
run:
	@echo "🚀 Running server..."
	@$(GO) run $(MAIN_FILE) -config $(CONFIG_FILE)

# ============================================================================
# 构建命令
# ============================================================================

## build: 构建生产二进制文件
build:
	@echo "📦 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

## build-migrate: 构建迁移工具
build-migrate:
	@echo "📦 Building migrate tool..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/migrate $(MIGRATE_FILE)
	@echo "✅ Build complete: $(BUILD_DIR)/migrate"

# ============================================================================
# 数据库命令
# ============================================================================

## migrate: 运行数据库迁移
migrate:
	@echo "🗄️  Running database migration..."
	@$(GO) run $(MIGRATE_FILE) -config $(CONFIG_FILE)

# ============================================================================
# 代码生成
# ============================================================================

## swagger: 生成 Swagger 文档
swagger:
	@echo "📝 Generating Swagger documentation..."
	@swag init -g $(MAIN_FILE) -o docs/swagger --parseDependency --parseInternal
	@echo "✅ Swagger docs generated"

## swagger-fmt: 格式化 Swagger 注释
swagger-fmt:
	@echo "📝 Formatting Swagger comments..."
	@swag fmt -g $(MAIN_FILE)

# ============================================================================
# 测试命令
# ============================================================================

## test: 运行所有测试
test:
	@echo "🧪 Running tests..."
	@$(GO) test -v ./...

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@$(GO) test -v -cover -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# ============================================================================
# 代码质量
# ============================================================================

## fmt: 格式化代码
fmt:
	@echo "🔧 Formatting code..."
	@$(GO) fmt ./...

## lint: 运行 linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

## tidy: 整理依赖
tidy:
	@echo "📦 Tidying dependencies..."
	@$(GO) mod tidy

# ============================================================================
# 清理命令
# ============================================================================

## clean: 清理构建产物
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

# ============================================================================
# 帮助
# ============================================================================

## help: 显示帮助信息
help:
	@echo ""
	@echo "LexVeritas Backend - Available Commands"
	@echo "========================================"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
