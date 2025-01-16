# Air 配置说明

Air 是一个用于 Go 项目的热重载工具，它可以监控项目文件的变化并自动重新编译和运行程序。

## 安装 Air

```bash
go install github.com/cosmtrek/air@latest
```

## 配置文件说明

项目根目录下的 `.air.toml` 配置文件详解：

```toml
# 工作目录
root = "."
tmp_dir = "tmp"

[build]
# 构建命令
cmd = "go build -o ./tmp/main ./cmd/qingshan-tools"
# 二进制文件路径
bin = "tmp/main"
# 需要监控的文件扩展名
include_ext = ["go", "tpl", "tmpl", "html"]
# 忽略的目录
exclude_dir = ["assets", "tmp", "vendor"]
# 文件变更后等待编译的时间（毫秒）
delay = 1000
# 终止信号等待时间
kill_delay = "0.5s"
# 构建日志文件
log = "build-errors.log"
# 是否发送中断信号
send_interrupt = false
# 发生错误时停止
stop_on_error = true

[color]
# 不同类型输出的颜色配置
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[log]
# 是否显示时间
time = false

[misc]
# 退出时清理临时文件
clean_on_exit = true
```

## 配置项说明

### 根配置
- `root`：项目根目录，通常设置为 "."
- `tmp_dir`：临时目录，用于存放编译后的二进制文件

### 构建配置 [build]
- `cmd`：构建命令，指定如何编译项目
- `bin`：编译后的二进制文件路径
- `include_ext`：需要监控的文件扩展名列表
- `exclude_dir`：不需要监控的目录列表
- `delay`：文件变更后等待编译的延迟时间（毫秒）
- `kill_delay`：发送终止信号后的等待时间
- `log`：构建错误日志文件
- `stop_on_error`：是否在构建错误时停止
- `send_interrupt`：是否发送中断信号而不是终止信号

### 颜色配置 [color]
- `main`：主要信息的颜色
- `watcher`：文件监控信息的颜色
- `build`：构建信息的颜色
- `runner`：运行信息的颜色

### 日志配置 [log]
- `time`：是否在日志中显示时间戳

### 其他配置 [misc]
- `clean_on_exit`：退出时是否清理临时文件

## 使用方法

1. 在项目根目录下创建 `.air.toml` 配置文件
2. 运行 `air` 命令启动热重载
3. 修改代码后，Air 会自动重新编译和运行程序

## 常用命令

```bash
# 启动热重载
air

# 使用指定配置文件
air -c .air.toml

# 查看版本
air -v
```

## 最佳实践

1. 合理设置 `exclude_dir`，排除不需要监控的目录
2. 适当调整 `delay` 时间，避免频繁编译
3. 配置 `log` 文件路径，方便排查构建错误
4. 开发环境建议设置 `stop_on_error = true`，及时发现编译错误 
