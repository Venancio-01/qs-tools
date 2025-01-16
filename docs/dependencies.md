# 项目依赖说明

本文档说明了项目中使用的主要依赖包及其用途。

## 主要依赖

### 1. github.com/spf13/cobra v1.8.1

- **用途**：命令行应用程序框架
- **主要功能**：
  - 提供现代化的 CLI 接口
  - 支持嵌套子命令
  - 自动生成命令行帮助信息
  - 支持命令行参数和标志
  - 智能建议功能

### 2. github.com/joho/godotenv v1.5.1

- **用途**：环境变量管理
- **主要功能**：
  - 从 .env 文件加载环境变量
  - 支持多环境配置
  - 简化配置管理

### 3. github.com/sirupsen/logrus v1.9.3

- **用途**：结构化日志记录
- **主要功能**：
  - 支持多种日志级别
  - 结构化日志输出
  - 支持多种输出格式（JSON、Text）
  - 支持日志字段定制
  - 高性能日志记录

## 间接依赖

### 1. github.com/spf13/pflag v1.0.5

- **用途**：命令行参数解析
- **说明**：Cobra 的依赖项，提供类 POSIX 风格的命令行标志处理

### 2. github.com/inconshreveable/mousetrap v1.1.0

- **用途**：Windows 命令行支持
- **说明**：用于在 Windows 环境下处理命令行应用程序的启动

### 3. golang.org/x/sys v0.29.0

- **用途**：系统调用接口
- **说明**：提供对底层操作系统功能的访问

### 4. 测试相关依赖

- github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
- github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2
- github.com/stretchr/testify v1.10.0
- **用途**：测试框架和工具
- **说明**：用于编写和运行测试用例，提供断言和测试辅助功能

## 版本管理

- 项目使用 Go Modules 进行依赖管理
- 所有依赖版本都在 go.mod 文件中明确指定
- 使用 go mod tidy 维护依赖的一致性

## 更新依赖

```bash
# 更新所有依赖到最新版本
go get -u ./...

# 更新特定依赖
go get -u github.com/spf13/cobra

# 整理依赖
go mod tidy
```
