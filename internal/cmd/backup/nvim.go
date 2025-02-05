package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"qs-tools/internal/utils"

	"github.com/spf13/cobra"
)

var nvimCmd = &cobra.Command{
	Use:   "nvim",
	Short: "备份 Neovim 配置",
	Long:  `备份 Neovim 配置文件并上传到远程服务器。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return backupNvim()
	},
}

func init() {
	BackupCmd.AddCommand(nvimCmd)
}

func backupNvim() error {
	fmt.Println("开始备份 Neovim 配置...")

	// 获取配置目录
	configDir := ""
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("LOCALAPPDATA"), "nvim")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户主目录失败: %v", err)
		}
		configDir = filepath.Join(homeDir, ".config", "nvim")
	}

	// 检查配置目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return fmt.Errorf("Neovim 配置目录不存在: %s", configDir)
	}

	// 创建临时目录
	tmpDir, cleanup, err := utils.CreateTempDir("nvim-backup")
	if err != nil {
		return err
	}
	defer cleanup()

	// 创建压缩文件
	backupFile := filepath.Join(tmpDir, "nvim_backup")
	if runtime.GOOS == "windows" {
		backupFile += ".zip"
	} else {
		backupFile += ".tar.gz"
	}

	if err := utils.CompressDir(configDir, backupFile); err != nil {
		return err
	}

	// 上传到远程服务器
	if err := utils.UploadToRemote("nvim", backupFile); err != nil {
		return err
	}

	fmt.Printf("\n✅ Neovim 配置备份成功！\n")
	return nil
}
