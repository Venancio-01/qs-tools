package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "qingshan-tools",
	Short: "青山工具集 - 一个实用的命令行工具集",
	Long: `青山工具集是一个集成了多种实用功能的命令行工具集。
可以帮助你完成各种日常任务，提高工作效率。`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
