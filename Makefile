.PHONY: all build clean run lint doc build-all

# 默认目标
all: lint build

# 构建当前平台版本
build:
	go build -v -o bin/qs-tools ./cmd/qs-tools

# 构建所有平台版本
build-all:
	./scripts/build.sh

# 代码检查
lint:
	golangci-lint run

# 生成 API 文档
doc:
	swag init -g cmd/qs-tools/main.go

# 运行开发服务器（热重载）
dev:
	air

# 清理构建文件
clean:
	rm -rf bin/ build/ tmp/
	go clean -cache

# 安装依赖工具
tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest 
