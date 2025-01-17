package backup

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Command 返回备份命令
func Command() *cobra.Command {
	return BackupCmd
}

// BackupCmd 表示备份命令
var BackupCmd = &cobra.Command{
	Use:   "backup [component]",
	Short: "备份配置文件",
	Long: `备份配置文件并上传到远程服务器。
目前支持的组件：
  - fish: 备份 Fish Shell 配置
  - scoop: 备份 Scoop 包管理器配置 (Windows)
  - nvim: 备份 Neovim 编辑器配置

支持的系统：
  - Ubuntu 及衍生版
  - Debian 及衍生版
  - Kylin (银河麒麟)
  - Windows`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请指定要备份的组件")
			return
		}

		var err error
		switch args[0] {
		case "fish":
			err = backupFish()
		case "scoop":
			err = backupScoop()
		case "nvim":
			err = backupNvim()
		default:
			fmt.Printf("不支持的组件: %s\n", args[0])
			return
		}

		if err != nil {
			fmt.Printf("备份失败: %v\n", err)
		}
	},
}
