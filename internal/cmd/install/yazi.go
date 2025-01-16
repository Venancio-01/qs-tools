package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func installYazi() error {
	fmt.Println("开始安装 Yazi 文件管理器...")

	// 1. 检查系统类型
	if runtime.GOOS == "windows" {
		return installYaziOnWindows()
	}

	// 2. 检查并安装依赖
	fmt.Println("1. 检查并安装必要的依赖...")
	deps := []string{
		"curl",
		"git",
		"gcc",
		"pkg-config",
		"make",
		"libmagic-dev",
		"build-essential",
		"libfontconfig-dev",
		"libglib2.0-dev",
	}
	installDepsCmd := exec.Command("sudo", append([]string{"apt", "install", "-y"}, deps...)...)
	installDepsCmd.Stdout = os.Stdout
	installDepsCmd.Stderr = os.Stderr
	if err := installDepsCmd.Run(); err != nil {
		return fmt.Errorf("安装依赖失败: %v", err)
	}

	// 3. 检查 Rust 工具链
	fmt.Println("\n2. 检查 Rust 工具链...")
	if _, err := exec.LookPath("cargo"); err != nil {
		fmt.Println("未检测到 Rust 工具链，正在安装...")
		installRustCmd := exec.Command("sh", "-c", "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y")
		installRustCmd.Stdout = os.Stdout
		installRustCmd.Stderr = os.Stderr
		if err := installRustCmd.Run(); err != nil {
			return fmt.Errorf("安装 Rust 失败: %v", err)
		}

		// 添加 Cargo 到 PATH
		cargoEnv := filepath.Join(os.Getenv("HOME"), ".cargo/env")
		if _, err := os.Stat(cargoEnv); err == nil {
			fmt.Printf("\n请运行以下命令加载 Rust 环境：\nsource %s\n\n然后重新运行安装命令。\n", cargoEnv)
			return nil
		}
	}

	// 更新 Rust 工具链
	fmt.Println("\n3. 更新 Rust 工具链...")
	updateCmd := exec.Command("rustup", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("更新 Rust 工具链失败: %v", err)
	}

	// 4. 克隆源码
	fmt.Println("\n4. 克隆源码...")
	tmpDir, err := os.MkdirTemp("", "yazi")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cloneCmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/sxyazi/yazi.git", tmpDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("克隆源码失败: %v", err)
	}

	// 5. 编译
	fmt.Println("\n5. 编译源码...")
	buildCmd := exec.Command("cargo", "build", "--release", "--locked")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("编译失败: %v", err)
	}

	// 6. 安装二进制文件
	fmt.Println("\n6. 安装二进制文件...")
	yaziPath := filepath.Join(tmpDir, "target/release/yazi")
	yaPath := filepath.Join(tmpDir, "target/release/ya")

	// 检查文件是否存在
	if _, err := os.Stat(yaziPath); err != nil {
		return fmt.Errorf("未找到 yazi 二进制文件: %v", err)
	}
	if _, err := os.Stat(yaPath); err != nil {
		return fmt.Errorf("未找到 ya 二进制文件: %v", err)
	}

	// 移动文件到 /usr/local/bin
	mvCmd := exec.Command("sudo", "mv", yaziPath, yaPath, "/usr/local/bin/")
	mvCmd.Stdout = os.Stdout
	mvCmd.Stderr = os.Stderr
	if err := mvCmd.Run(); err != nil {
		return fmt.Errorf("安装二进制文件失败: %v", err)
	}

	fmt.Println("\n✅ Yazi 安装成功！")
	printYaziUsage()
	return nil
}

func installYaziOnWindows() error {
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

	// 安装 yazi
	fmt.Println("安装 Yazi...")
	cmd = exec.Command("scoop", "install", "yazi")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("安装 Yazi 失败: %v", err)
	}

	fmt.Println("\n✅ Yazi 安装成功！")
	printYaziUsage()
	return nil
}

func printYaziUsage() {
	fmt.Println("\n使用说明：")
	fmt.Println("1. 在终端中输入 'yazi' 启动文件管理器")
	fmt.Println("   或者使用 'ya' 快捷命令")
	fmt.Println("\n2. 常用快捷键：")
	fmt.Println("   - h/j/k/l: 导航")
	fmt.Println("   - Space: 预览文件")
	fmt.Println("   - Enter: 打开文件/目录")
	fmt.Println("   - y: 复制")
	fmt.Println("   - d: 剪切")
	fmt.Println("   - p: 粘贴")
	fmt.Println("   - q: 退出")
}
