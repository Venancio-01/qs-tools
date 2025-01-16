package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func installAsdf() {
	// 检查是否为支持的系统
	if !isDebianBased() {
		fmt.Println("当前系统不是基于 Debian 的系统（如 Ubuntu、Debian、Kylin 等）")
		return
	}

	fmt.Println("开始安装 asdf 版本管理器...")

	// 安装依赖
	fmt.Println("\n1. 安装必要的依赖...")
	deps := []string{
		"curl",
		"git",
		"make",
		"unzip",
	}

	installDepsCmd := exec.Command("sudo", append([]string{"apt", "install", "-y"}, deps...)...)
	installDepsCmd.Stdout = os.Stdout
	installDepsCmd.Stderr = os.Stderr
	if err := installDepsCmd.Run(); err != nil {
		fmt.Printf("安装依赖失败: %v\n", err)
		return
	}

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return
	}

	// 检查是否已安装
	asdfDir := filepath.Join(homeDir, ".asdf")
	if _, err := os.Stat(asdfDir); err == nil {
		fmt.Println("\nasdf 已经安装。如需重新安装，请先删除 ~/.asdf 目录")
		return
	}

	// 克隆 asdf 仓库
	fmt.Println("\n2. 克隆 asdf 仓库...")
	cloneCmd := exec.Command("git", "clone", "https://github.com/asdf-vm/asdf.git", asdfDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("克隆 asdf 仓库失败: %v\n", err)
		return
	}

	// 检查当前 Shell
	currentShell := os.Getenv("SHELL")
	if currentShell == "" {
		currentShell = "/bin/bash" // 默认使用 bash
	}

	// 配置 Shell 集成
	fmt.Println("\n3. 配置 Shell 集成...")
	shellConfigFile := ""
	shellInitCmd := ""

	switch {
	case strings.Contains(currentShell, "fish"):
		shellConfigFile = filepath.Join(homeDir, ".config/fish/config.fish")
		shellInitCmd = "source ~/.asdf/asdf.fish"
	case strings.Contains(currentShell, "zsh"):
		shellConfigFile = filepath.Join(homeDir, ".zshrc")
		shellInitCmd = `. $HOME/.asdf/asdf.sh`
	default: // bash
		shellConfigFile = filepath.Join(homeDir, ".bashrc")
		shellInitCmd = `. $HOME/.asdf/asdf.sh`
	}

	// 确保配置目录存在
	if err := os.MkdirAll(filepath.Dir(shellConfigFile), 0755); err != nil {
		fmt.Printf("创建配置目录失败: %v\n", err)
		return
	}

	// 添加初始化命令到 Shell 配置
	f, err := os.OpenFile(shellConfigFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("打开 Shell 配置文件失败: %v\n", err)
		return
	}
	defer f.Close()

	// 检查是否已经配置
	content, err := os.ReadFile(shellConfigFile)
	if err == nil && !strings.Contains(string(content), "asdf.") {
		if _, err := f.WriteString("\n# asdf 版本管理器\n" + shellInitCmd + "\n"); err != nil {
			fmt.Printf("写入 Shell 配置失败: %v\n", err)
			return
		}
	}

	fmt.Println("\n✅ asdf 安装成功！")
	fmt.Println("\n使用说明：")
	fmt.Println("1. 重新打开终端或执行以下命令使配置生效：")
	fmt.Printf("   source %s\n", shellConfigFile)
	fmt.Println("\n2. 常用命令：")
	fmt.Println("   - 查看可用插件：asdf plugin list all")
	fmt.Println("   - 安装插件：asdf plugin add <name>")
	fmt.Println("   - 安装特定版本：asdf install <name> <version>")
	fmt.Println("   - 设置全局版本：asdf global <name> <version>")
	fmt.Println("   - 设置本地版本：asdf local <name> <version>")
	fmt.Println("   - 查看当前版本：asdf current")
	fmt.Println("\n3. 示例 - 安装 Node.js：")
	fmt.Println("   asdf plugin add nodejs")
	fmt.Println("   asdf install nodejs latest")
	fmt.Println("   asdf global nodejs latest")
} 
