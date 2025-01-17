package apply

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
	Short: "恢复 Neovim 配置",
	Long:  `从远程服务器下载并恢复 Neovim 的配置文件。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return applyNvim()
	},
}

func init() {
	ApplyCmd.AddCommand(nvimCmd)
}

func applyNvim() error {
	fmt.Println("开始恢复 Neovim 配置...")

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

	// 创建临时目录
	tmpDir, cleanup, err := utils.CreateTempDir("nvim-restore")
	if err != nil {
		return err
	}
	defer cleanup()

	// 从远程服务器下载备份文件
	backupFile := filepath.Join(tmpDir, "nvim_backup")
	if runtime.GOOS == "windows" {
		backupFile += ".zip"
	} else {
		backupFile += ".tar.gz"
	}

	if err := utils.DownloadFromRemote("nvim", backupFile); err != nil {
		return err
	}

	// 解压配置文件
	if err := utils.ExtractFile(backupFile, filepath.Dir(configDir)); err != nil {
		return err
	}

	fmt.Printf("\n✅ Neovim 配置恢复成功！\n")
	return nil
}
