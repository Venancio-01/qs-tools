package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示程序版本",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("QS工具集 %s\n", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
