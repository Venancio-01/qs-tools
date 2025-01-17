package apply

import (
	"fmt"
	"github.com/spf13/cobra"
)

// Command 返回恢复命令
func Command() *cobra.Command {
	return ApplyCmd
}

// ApplyCmd 表示 apply 命令
var ApplyCmd = &cobra.Command{
	Use:   "apply [component]",
	Short: "从远程服务器恢复配置",
	Long: `从远程服务器下载并恢复之前备份的配置文件。
目前支持的组件：
  - fish: 恢复 Fish Shell 配置
  - scoop: 恢复 Scoop 包管理器配置 (Windows)
  - nvim: 恢复 Neovim 编辑器配置

支持的系统：
  - Ubuntu 及衍生版
  - Debian 及衍生版
  - Kylin (银河麒麟)
  - Windows`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请指定要恢复的组件")
			return
		}

		var err error
		switch args[0] {
		case "fish":
			err = applyFish()
		case "scoop":
			err = applyScoop()
		case "nvim":
			err = applyNvim()
		default:
			fmt.Printf("不支持的组件: %s\n", args[0])
			return
		}

		if err != nil {
			fmt.Printf("恢复失败: %v\n", err)
		}
	},
}
