package install

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Command 返回安装命令
func Command() *cobra.Command {
	return Cmd
}

var Cmd = &cobra.Command{
	Use:   "install [component]",
	Short: "安装常用工具和软件",
	Long: `安装常用的工具和软件。
目前支持的组件：
  - fish: 安装 Fish Shell
  - asdf: 安装 asdf 版本管理器
  - yazi: 安装 Yazi 文件管理器
  - scoop: 安装 Scoop 包管理器 (Windows)
  - nvim: 安装 Neovim 编辑器

支持的系统：
  - Ubuntu 及衍生版
  - Debian 及衍生版
  - Kylin (银河麒麟)
  - Windows`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请指定要安装的组件")
			return
		}

		var err error
		switch args[0] {
		case "fish":
			installFish()
		case "asdf":
			installAsdf()
		case "yazi":
			err = installYazi()
		case "scoop":
			installScoop()
		case "nvim":
			err = installNvim()
		default:
			fmt.Printf("不支持的组件: %s\n", args[0])
			return
		}

		if err != nil {
			fmt.Printf("安装失败: %v\n", err)
		}
	},
}
