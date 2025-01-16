# qs-tools

个人工具箱，用于管理和同步各种开发工具的配置。

## 功能特性

1. 安装常用工具和软件
   - Fish Shell
   - asdf 版本管理器
   - Yazi 文件管理器
   - Scoop 包管理器 (Windows)

2. 备份和恢复配置
   - Fish Shell 配置
   - Scoop 包管理器配置 (Windows)

## 系统要求

- Linux (Ubuntu/Debian/Kylin) 或 Windows
- Go 1.20 或更高版本

## 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/qs-tools.git
cd qs-tools

# 编译
go build -o qs-tools

# 将二进制文件移动到 PATH 目录（可选）
sudo mv qs-tools /usr/local/bin/
```

## 使用方法

### 安装工具

```bash
# 安装 Fish Shell
qs-tools install fish

# 安装 asdf 版本管理器
qs-tools install asdf

# 安装 Yazi 文件管理器
qs-tools install yazi

# 安装 Scoop 包管理器 (Windows)
qs-tools install scoop
```

### 备份配置

```bash
# 备份 Fish Shell 配置
qs-tools backup fish

# 备份 Scoop 配置 (Windows)
qs-tools backup scoop
```

### 恢复配置

```bash
# 恢复 Fish Shell 配置
qs-tools apply fish
```

## 支持的系统

- Ubuntu 及衍生版
- Debian 及衍生版
- Kylin (银河麒麟)
- Windows

## 配置说明

1. Fish Shell
   - 备份 `~/.config/fish` 目录下的所有配置文件
   - 自动上传到远程服务器
   - 支持一键恢复

2. Scoop (Windows)
   - 备份已安装的应用列表
   - 备份软件源列表
   - 备份配置文件
   - 生成恢复脚本
   - 自动上传到远程服务器

## 注意事项

1. 备份功能需要网络连接
2. Windows 下需要安装 OpenSSH 客户端
3. 恢复配置前会自动备份现有配置
4. 某些功能可能需要管理员/超级用户权限

## 开发计划

- [ ] 添加更多工具的支持
- [ ] 添加配置文件的版本控制
- [ ] 添加配置文件的差异比较
- [ ] 添加 Web 界面
- [ ] 添加自动更新功能

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件
