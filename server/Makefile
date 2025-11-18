.PHONY: help build run dev test clean docker-up docker-down migrate

# 默认目标
help:
	@echo "Chatty Server - 可用命令:"
	@echo "  make build       - 编译服务器"
	@echo "  make run         - 运行服务器"
	@echo "  make dev         - 开发模式(热重载)"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理构建文件"
	@echo "  make docker-up   - 启动Docker服务(PostgreSQL & Redis)"
	@echo "  make docker-down - 停止Docker服务"
	@echo "  make migrate     - 运行数据库迁移"
	@echo "  make lint        - 运行代码检查"

# 编译
build:
	@echo "编译服务器..."
	@go build -o bin/chatty-server cmd/server/main.go
	@echo "编译完成! 二进制文件: bin/chatty-server"

# 运行
run: build
	@echo "启动服务器..."
	@./bin/chatty-server

# 开发模式(需要安装air: go install github.com/air-verse/air@latest)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air未安装. 安装命令: go install github.com/air-verse/air@latest"; \
		echo "使用普通模式运行..."; \
		go run cmd/server/main.go; \
	fi

# 测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 测试覆盖率
test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 清理
clean:
	@echo "清理构建文件..."
	@rm -f bin/chatty-server
	@rm -f coverage.out coverage.html
	@echo "清理完成"

# Docker服务
docker-up:
	@echo "启动Docker服务..."
	@docker-compose up -d
	@echo "等待服务就绪..."
	@sleep 3
	@docker-compose ps

docker-down:
	@echo "停止Docker服务..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

# 数据库迁移
migrate:
	@echo "运行数据库迁移..."
	@go run cmd/server/main.go migrate

# 代码检查
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint未安装"; \
		echo "安装: https://golangci-lint.run/welcome/install/"; \
	fi

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@echo "格式化完成"

# 安装开发工具
install-tools:
	@echo "安装开发工具..."
	@go install github.com/air-verse/air@latest
	@echo "开发工具安装完成"

# 依赖管理
deps:
	@echo "下载依赖..."
	@go mod download
	@go mod tidy
	@echo "依赖更新完成"
