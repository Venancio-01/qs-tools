package install

import (
	"fmt"
	"os"
	"os/exec"
)

// InstallFish 安装 Fish Shell
func installFish() {
	// 检查是否为支持的系统
	if !isDebianBased() {
		fmt.Println("当前系统不是基于 Debian 的系统（如 Ubuntu、Debian、Kylin 等）")
		return
	}

	fmt.Println("检测到支持的系统，开始安装 Fish Shell...")

	// 如果是麒麟系统，使用直接下载安装的方式
	if isKylin() {
		fmt.Println("\n检测到 Kylin 系统，将使用直接下载安装方式...")

		// 创建临时目录
		tmpDir := "/tmp/qs-tools-fish"
		os.MkdirAll(tmpDir, 0755)
		debPath := tmpDir + "/fish.deb"

		// 下载 deb 包 (使用 3.7.1 版本)
		fmt.Println("\n1. 下载 Fish Shell 安装包...")
		downloadCmd := exec.Command("wget",
			"https://download.opensuse.org/repositories/shells:/fish:/release:/3/Debian_10/amd64/fish_3.7.1-1_amd64.deb",
			"-O", debPath)
		downloadCmd.Stdout = os.Stdout
		downloadCmd.Stderr = os.Stderr
		if err := downloadCmd.Run(); err != nil {
			fmt.Printf("下载安装包失败: %v\n", err)
			fmt.Println("\n尝试安装系统默认版本...")
			// 如果下载失败，尝试使用系统默认源安装
			installFromApt()
			return
		}

		// 安装 deb 包前先安装依赖
		fmt.Println("\n2. 安装必要的依赖...")
		depsCmd := exec.Command("sudo", "apt-get", "install", "-y", "libpcre2-32-0")
		depsCmd.Stdout = os.Stdout
		depsCmd.Stderr = os.Stderr
		depsCmd.Run() // 忽略错误，让后续的 apt-get install -f 处理

		// 安装 deb 包
		fmt.Println("\n3. 安装 Fish Shell...")
		installCmd := exec.Command("sudo", "dpkg", "-i", debPath)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			fmt.Println("\n安装过程中可能缺少依赖，尝试修复...")
			fixCmd := exec.Command("sudo", "apt-get", "install", "-f", "-y")
			fixCmd.Stdout = os.Stdout
			fixCmd.Stderr = os.Stderr
			if err := fixCmd.Run(); err != nil {
				fmt.Printf("修复依赖失败: %v\n", err)
				fmt.Println("\n尝试安装系统默认版本...")
				// 如果安装失败，尝试使用系统默认源安装
				installFromApt()
				return
			}

			// 重试安装
			if err := installCmd.Run(); err != nil {
				fmt.Printf("安装失败: %v\n", err)
				fmt.Println("\n尝试安装系统默认版本...")
				// 如果安装失败，尝试使用系统默认源安装
				installFromApt()
				return
			}
		}

		// 清理临时文件
		os.RemoveAll(tmpDir)
	} else {
		installFromApt()
	}

	fmt.Println("\n✅ Fish Shell 安装成功！")

	// 获取 fish 版本
	versionCmd := exec.Command("fish", "--version")
	versionOutput, err := versionCmd.Output()
	if err == nil {
		fmt.Printf("\n当前安装的版本：%s", string(versionOutput))
	}

	fmt.Println("\n你可以通过以下命令将 Fish 设置为默认 Shell：")
	fmt.Println("chsh -s $(which fish)")

	// 检查是否为 Kylin 系统，提供额外说明
	if isKylin() {
		fmt.Println("\n注意：在 Kylin 系统上，你可能需要以下额外步骤：")
		fmt.Println("1. 编辑 /etc/shells 文件，添加 Fish Shell 路径")
		fmt.Println("   sudo echo $(which fish) >> /etc/shells")
		fmt.Println("2. 然后再执行更改默认 Shell 的命令")
	}

	fmt.Println("\n💡 提示：首次启动 Fish Shell 时，建议运行以下命令完成初始配置：")
	fmt.Println("fish_config")
}

// installFromApt 使用 apt 安装 Fish Shell
func installFromApt() {
	// 其他 Debian 系统使用 PPA 安装
	needPPA := true

	// 检查并安装必要的依赖
	fmt.Println("\n1. 检查并安装必要的依赖...")
	checkDepsCmd := exec.Command("which", "apt-add-repository")
	if err := checkDepsCmd.Run(); err != nil {
		fmt.Println("正在安装 software-properties-common...")
		installDepsCmd := exec.Command("sudo", "apt", "install", "-y", "software-properties-common")
		installDepsCmd.Stdout = os.Stdout
		installDepsCmd.Stderr = os.Stderr
		if err := installDepsCmd.Run(); err != nil {
			fmt.Printf("安装依赖失败: %v\n", err)
			fmt.Println("将尝试使用系统默认软件源安装 Fish Shell...")
			needPPA = false
		}
	}

	// 如果需要且可以添加 PPA，则添加
	if needPPA {
		fmt.Println("\n2. 添加 Fish Shell 官方 PPA...")
		addRepoCmd := exec.Command("sudo", "apt-add-repository", "-y", "ppa:fish-shell/release-3")
		addRepoCmd.Stdout = os.Stdout
		addRepoCmd.Stderr = os.Stderr
		if err := addRepoCmd.Run(); err != nil {
			fmt.Printf("添加 Fish Shell PPA 失败: %v\n", err)
			fmt.Println("将尝试使用系统默认软件源安装...")
		}
	}

	// 更新包索引
	fmt.Println("\n3. 更新软件包索引...")
	updateCmd := exec.Command("sudo", "apt", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		fmt.Printf("更新包索引失败: %v\n", err)
		return
	}

	// 安装 fish
	fmt.Println("\n4. 安装 Fish Shell...")
	installCmd := exec.Command("sudo", "apt", "install", "-y", "fish")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Printf("安装 Fish Shell 失败: %v\n", err)
		return
	}
} 
