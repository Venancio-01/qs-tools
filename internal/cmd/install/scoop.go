package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func installScoop() {
	if runtime.GOOS != "windows" {
		fmt.Println("Scoop 仅支持 Windows 系统")
		return
	}

	fmt.Println("开始安装 Scoop 包管理器...")

	// 检查是否已安装
	if _, err := exec.LookPath("scoop"); err == nil {
		fmt.Println("\nScoop 已经安装。")
		printScoopUsage()
		return
	}

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return
	}

	// 设置 SCOOP 环境变量
	scoopDir := filepath.Join(homeDir, "scoop")
	os.Setenv("SCOOP", scoopDir)

	// 使用 PowerShell 安装 Scoop
	fmt.Println("\n1. 下载并安装 Scoop...")
	installCmd := exec.Command("powershell", "-Command",
		`Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://get.scoop.sh'))`)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Printf("安装 Scoop 失败: %v\n", err)
		return
	}

	// 添加常用 bucket
	fmt.Println("\n2. 添加常用软件源...")
	buckets := []string{"extras", "versions", "nerd-fonts", "java"}
	for _, bucket := range buckets {
		fmt.Printf("添加 %s bucket...\n", bucket)
		addCmd := exec.Command("scoop", "bucket", "add", bucket)
		addCmd.Stdout = os.Stdout
		addCmd.Stderr = os.Stderr
		addCmd.Run() // 忽略错误，因为可能已经添加
	}

	fmt.Println("\n✅ Scoop 安装成功！")
	printScoopUsage()
}

func printScoopUsage() {
	fmt.Println("\n使用说明：")
	fmt.Println("1. 基本命令：")
	fmt.Println("   - 搜索软件：scoop search <app>")
	fmt.Println("   - 安装软件：scoop install <app>")
	fmt.Println("   - 更新软件：scoop update <app>")
	fmt.Println("   - 卸载软件：scoop uninstall <app>")
	fmt.Println("   - 查看已安装：scoop list")
	fmt.Println("   - 清理缓存：scoop cleanup")
	fmt.Println("\n2. 软件源管理：")
	fmt.Println("   - 添加源：scoop bucket add <bucket>")
	fmt.Println("   - 移除源：scoop bucket rm <bucket>")
	fmt.Println("   - 查看已添加源：scoop bucket list")
	fmt.Println("\n3. 系统维护：")
	fmt.Println("   - 更新 Scoop：scoop update")
	fmt.Println("   - 更新所有软件：scoop update *")
	fmt.Println("   - 检查问题：scoop checkup")
	fmt.Println("\n4. 已添加的软件源：")
	fmt.Println("   - extras: 包含大量常用软件")
	fmt.Println("   - versions: 包含软件的多个版本")
	fmt.Println("   - nerd-fonts: 包含编程字体")
	fmt.Println("   - java: 包含 Java 相关软件")
	fmt.Println("\n5. 推荐安装的基础软件：")
	fmt.Println("   scoop install git 7zip sudo curl")
} 
