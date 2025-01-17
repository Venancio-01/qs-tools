package apply

import (
	"fmt"
	"os"
	"path/filepath"

	"qs-tools/internal/utils"

	"github.com/spf13/cobra"
)

var fishCmd = &cobra.Command{
	Use:   "fish",
	Short: "恢复 Fish Shell 配置",
	Long:  `从远程服务器下载并恢复 Fish Shell 的配置文件。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return applyFish()
	},
}

func init() {
	ApplyCmd.AddCommand(fishCmd)
}

func applyFish() error {
	fmt.Println("开始恢复 Fish Shell 配置...")

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}

	// 创建临时目录
	tmpDir, cleanup, err := utils.CreateTempDir("fish-restore")
	if err != nil {
		return err
	}
	defer cleanup()

	// 从远程服务器下载备份文件
	backupFile := filepath.Join(tmpDir, "fish_backup.tar.gz")
	if err := utils.DownloadFromRemote("fish", backupFile); err != nil {
		return err
	}

	// 解压配置文件
	configDir := filepath.Join(homeDir, ".config")
	if err := utils.ExtractFile(backupFile, configDir); err != nil {
		return err
	}

	fmt.Printf("\n✅ Fish Shell 配置恢复成功！\n")
	return nil
}
