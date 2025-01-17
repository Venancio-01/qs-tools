package backup

import (
	"fmt"
	"os"
	"path/filepath"

	"qs-tools/internal/utils"

	"github.com/spf13/cobra"
)

var fishCmd = &cobra.Command{
	Use:   "fish",
	Short: "备份 Fish Shell 配置",
	Long:  `备份 Fish Shell 配置文件并上传到远程服务器。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return backupFish()
	},
}

func init() {
	BackupCmd.AddCommand(fishCmd)
}

func backupFish() error {
	fmt.Println("开始备份 Fish Shell 配置...")

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}

	// Fish 配置目录
	configDir := filepath.Join(homeDir, ".config", "fish")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("Fish 配置目录不存在: %s", configDir)
	}

	// 创建临时目录
	tmpDir, cleanup, err := utils.CreateTempDir("fish-backup")
	if err != nil {
		return err
	}
	defer cleanup()

	// 创建压缩文件
	backupFile := filepath.Join(tmpDir, "fish_backup.tar.gz")
	if err := utils.CompressDir(configDir, backupFile); err != nil {
		return err
	}

	// 上传到远程服务器
	if err := utils.UploadToRemote("fish", backupFile); err != nil {
		return err
	}

	fmt.Printf("\n✅ Fish Shell 配置备份成功！\n")
	return nil
}
