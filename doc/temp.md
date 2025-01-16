# 现代 Go 项目最佳实践指南

## 1. 项目结构

```
project_name/
├── .vscode/                    # VSCode 配置
│   ├── launch.json            # 调试配置
│   ├── tasks.json            # 任务配置
│   └── settings.json         # 编辑器设置
├── cmd/                      # 主要应用入口
│   └── project_name/        # 主程序
│       └── main.go
├── internal/                # 私有应用和库代码
│   ├── api/                # API 接口
│   ├── config/             # 配置
│   ├── middleware/         # 中间件
│   ├── model/             # 数据模型
│   ├── repository/        # 数据仓储
│   ├── service/           # 业务逻辑
│   └── utils/             # 工具函数
├── pkg/                    # 可被外部使用的库代码
├── api/                    # OpenAPI/Swagger 规范
│   └── swagger.yaml
├── web/                    # Web 静态资源
├── configs/                # 配置文件目录
│   ├── config.yaml
│   └── config.dev.yaml
├── deployments/            # 部署配置和模板
│   ├── docker/
│   └── kubernetes/
├── docs/                   # 文档
├── scripts/               # 构建和部署脚本
├── test/                  # 测试目录
├── .air.toml             # Air 配置（热重载）
├── .env                   # 环境变量
├── .gitignore
├── Dockerfile
├── go.mod                # Go 模块文件
├── go.sum                # Go 依赖校验
├── Makefile             # 项目管理命令
└── README.md
```

## 2. 开发环境配置

### 2.1 必要工具

- Go 1.21+
- VSCode
- Git
- Docker（可选）
- Make
- Air（热重载）
- golangci-lint（代码检查）
- swag（API 文档）
- mockgen（测试模拟）

### 2.2 VSCode 扩展

- Go
- Go Test Explorer
- Go Coverage
- Go Outliner
- GitLens
- Docker
- YAML
- REST Client

### 2.3 关键配置文件

#### .vscode/settings.json

```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintFlags": ["--fast"],
    "go.formatTool": "goimports",
    "go.testOnSave": true,
    "go.coverOnSave": true,
    "editor.formatOnSave": true,
    "[go]": {
        "editor.defaultFormatter": "golang.go"
    }
}
```

#### .golangci.yml

```yaml
linters:
  enable:
    - gofmt
    - golint
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - misspell

run:
  deadline: 5m

issues:
  exclude-use-default: false
```

#### Makefile

```makefile
.PHONY: all build test clean run lint doc

# 默认目标
all: lint test build

# 构建
build:
 go build -v -o bin/app ./cmd/project_name

# 运行测试
test:
 go test -v -race -cover ./...

# 代码检查
lint:
 golangci-lint run

# 生成 API 文档
doc:
 swag init -g cmd/project_name/main.go

# 运行开发服务器（热重载）
dev:
 air

# 清理构建文件
clean:
 rm -rf bin/
 go clean -cache

# 安装依赖工具
tools:
 go install github.com/cosmtrek/air@latest
 go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
 go install github.com/swaggo/swag/cmd/swag@latest
 go install github.com/golang/mock/mockgen@latest
```

## 3. 开发最佳实践

### 3.1 代码规范

- 遵循官方 Go 代码规范
- 使用 gofmt/goimports 格式化代码
- 实现接口优先设计
- 使用依赖注入
- 错误处理遵循 Go 1.13+ 错误链
- 合理使用泛型（Go 1.18+）

### 3.2 项目布局

- 遵循 Standard Go Project Layout
- 使用 internal 封装私有代码
- 使用 pkg 存放可复用库
- 使用 cmd 存放入口程序

### 3.3 依赖管理

```bash
# 初始化模块
go mod init project_name

# 添加依赖
go get -u package_name

# 整理依赖
go mod tidy

# 验证依赖
go mod verify

# 更新所有依赖
go get -u all
```

### 3.4 测试规范

- 单元测试（_test.go）
- 基准测试（Benchmark）
- 示例测试（Example）
- 使用 testify 断言库
- 使用 mockgen 生成 mock
- 保持测试覆盖率 > 80%

## 4. 部署流程

### 4.1 容器化部署

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache make git

# 复制源代码
COPY . .

# 构建应用
RUN make build

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/bin/app .
COPY --from=builder /app/configs ./configs

# 设置环境变量
ENV GO_ENV=production

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./app"]
```

### 4.2 CI/CD 配置

- 使用 GitHub Actions 或 GitLab CI
- 自动化测试和构建
- 代码质量检查
- 自动发布版本
- 容器镜像构建

### 4.3 监控和日志

- 使用 zap 或 zerolog 记录日志
- 集成 Prometheus 指标
- 使用 Jaeger 进行链路追踪
- 健康检查接口
- 优雅关闭处理

## 5. 开发工作流

### 5.1 日常开发

```bash
# 启动开发服务器（热重载）
make dev

# 运行测试
make test

# 代码检查
make lint

# 生成 API 文档
make doc

# 构建应用
make build
```

### 5.2 版本发布

1. 更新版本号
2. 更新 CHANGELOG
3. 运行完整测试套件
4. 创建 Git Tag
5. 触发 CI/CD 流程
6. 部署到生产环境

### 5.3 性能优化

- 使用 pprof 进行性能分析
- 合理使用 goroutine 和通道
- 注意内存分配和垃圾回收
- 使用连接池管理资源
- 实现合理的缓存策略
