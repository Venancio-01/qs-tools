package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [component]",
	Short: "安装常用工具和软件",
	Long: `安装常用的工具和软件。
目前支持的组件：
  - fish: 安装 Fish Shell`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请指定要安装的组件")
			return
		}

		switch args[0] {
		case "fish":
			installFish()
		default:
			fmt.Printf("不支持的组件: %s\n", args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}

func installFish() {
	// 检查是否为类 Ubuntu 系统
	if !isUbuntuLike() {
		fmt.Println("当前系统不是 Ubuntu 或类 Ubuntu 系统")
		return
	}

	fmt.Println("检测到 Ubuntu 或类 Ubuntu 系统，开始安装 Fish Shell...")

	// 更新包索引
	updateCmd := exec.Command("sudo", "apt", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		fmt.Printf("更新包索引失败: %v\n", err)
		return
	}

	// 安装 fish
	installCmd := exec.Command("sudo", "apt", "install", "-y", "fish")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		fmt.Printf("安装 Fish Shell 失败: %v\n", err)
		return
	}

	fmt.Println("Fish Shell 安装成功！")
	fmt.Println("你可以通过以下命令将 Fish 设置为默认 Shell：")
	fmt.Println("chsh -s $(which fish)")
}

func isUbuntuLike() bool {
	// 检查 /etc/os-release 文件
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}

	content := string(data)
	return strings.Contains(content, "Ubuntu") ||
		strings.Contains(content, "Debian") ||
		strings.Contains(content, "LinuxMint") ||
		strings.Contains(content, "Pop!_OS")
}
