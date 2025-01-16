package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func installNvim() error {
	fmt.Println("开始安装 Neovim...")

	if runtime.GOOS == "windows" {
		return installNvimOnWindows()
	}

	// 1. 检查并安装依赖
	fmt.Println("1. 检查并安装必要的依赖...")
	deps := []string{
		"ninja-build",
		"gettext",
		"cmake",
		"unzip",
		"curl",
		"git",
	}
	installDepsCmd := exec.Command("sudo", append([]string{"apt", "install", "-y"}, deps...)...)
	installDepsCmd.Stdout = os.Stdout
	installDepsCmd.Stderr = os.Stderr
	if err := installDepsCmd.Run(); err != nil {
		return fmt.Errorf("安装依赖失败: %v", err)
	}

	// 2. 创建临时目录
	fmt.Println("\n2. 创建临时目录...")
	tmpDir, err := os.MkdirTemp("", "nvim")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 3. 克隆源码
	fmt.Println("3. 克隆 Neovim 源码...")
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/neovim/neovim.git", tmpDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("克隆源码失败: %v", err)
	}

	// 4. 编译安装
	fmt.Println("\n4. 编译安装...")
	buildCmd := exec.Command("make", "CMAKE_BUILD_TYPE=Release")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("编译失败: %v", err)
	}

	// 5. 安装
	fmt.Println("\n5. 安装到系统...")
	installCmd := exec.Command("sudo", "make", "install")
	installCmd.Dir = tmpDir
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("安装失败: %v", err)
	}

	// 6. 安装配置管理器（可选）
	if err := installNvimConfig(); err != nil {
		fmt.Printf("\n⚠️ 配置安装失败: %v\n", err)
	}

	fmt.Println("\n✅ Neovim 安装成功！")
	printNvimUsage()
	return nil
}

func installNvimOnWindows() error {
	// 检查是否安装了 scoop
	if _, err := exec.LookPath("scoop"); err != nil {
		fmt.Println("未检测到 Scoop，正在安装...")
		cmd := exec.Command("powershell", "-Command", "Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force; iwr -useb get.scoop.sh | iex")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("安装 Scoop 失败: %v", err)
		}
	}

	// 添加 extras bucket
	fmt.Println("添加 extras bucket...")
	cmd := exec.Command("scoop", "bucket", "add", "extras")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run() // 忽略错误，因为可能已经添加

	// 安装 neovim
	fmt.Println("安装 Neovim...")
	cmd = exec.Command("scoop", "install", "neovim")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("安装 Neovim 失败: %v", err)
	}

	// 安装配置管理器（可选）
	if err := installNvimConfig(); err != nil {
		fmt.Printf("\n⚠️ 配置安装失败: %v\n", err)
	}

	fmt.Println("\n✅ Neovim 安装成功！")
	printNvimUsage()
	return nil
}

func installNvimConfig() error {
	// 获取配置目录
	configDir := ""
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "nvim")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %v", err)
		}
		configDir = filepath.Join(homeDir, ".config", "nvim")
	}

	// 备份现有配置
	if _, err := os.Stat(configDir); err == nil {
		backupDir := configDir + ".bak"
		if err := os.RemoveAll(backupDir); err != nil {
			return fmt.Errorf("删除旧的备份失败: %v", err)
		}
		if err := os.Rename(configDir, backupDir); err != nil {
			return fmt.Errorf("备份现有配置失败: %v", err)
		}
		fmt.Printf("已备份现有配置到: %s\n", backupDir)
	}

	// 克隆配置仓库
	fmt.Println("\n6. 安装配置文件...")
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/LazyVim/starter", configDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("克隆配置仓库失败: %v", err)
	}

	// 删除 .git 目录
	if err := os.RemoveAll(filepath.Join(configDir, ".git")); err != nil {
		return fmt.Errorf("删除 .git 目录失败: %v", err)
	}

	return nil
}

func printNvimUsage() {
	fmt.Println("\n使用说明：")
	fmt.Println("1. 在终端中输入 'nvim' 启动编辑器")
	fmt.Println("\n2. 首次启动会自动安装插件，请耐心等待")
	fmt.Println("\n3. 常用快捷键：")
	fmt.Println("   - Space: 打开命令面板")
	fmt.Println("   - Space + e: 打开文件浏览器")
	fmt.Println("   - Space + ff: 查找文件")
	fmt.Println("   - Space + fg: 全局搜索")
	fmt.Println("   - Space + qq: 退出")
	fmt.Println("\n4. 如果需要恢复原有配置，可以删除配置目录后还原备份")
	if runtime.GOOS == "windows" {
		fmt.Printf("   配置目录: %s\n", filepath.Join(os.Getenv("LOCALAPPDATA"), "nvim"))
	} else {
		fmt.Printf("   配置目录: ~/.config/nvim\n")
	}
}
