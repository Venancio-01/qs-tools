package apply

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"qs-tools/internal/utils"

	"github.com/spf13/cobra"
)

var scoopCmd = &cobra.Command{
	Use:   "scoop",
	Short: "恢复 Scoop 配置",
	Long:  `从远程服务器下载并恢复 Scoop 的配置文件。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return applyScoop()
	},
}

func init() {
	ApplyCmd.AddCommand(scoopCmd)
}

func applyScoop() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Scoop 仅支持 Windows 系统")
	}

	fmt.Println("开始恢复 Scoop 配置...")

	// 创建临时目录
	tmpDir, cleanup, err := utils.CreateTempDir("scoop-restore")
	if err != nil {
		return err
	}
	defer cleanup()

	// 从远程服务器下载备份文件
	backupFile := filepath.Join(tmpDir, "scoop_backup.zip")
	if err := utils.DownloadFromRemote("scoop", backupFile); err != nil {
		return err
	}

	// 解压配置文件
	if err := utils.ExtractFile(backupFile, tmpDir); err != nil {
		return err
	}

	// 执行恢复脚本
	fmt.Println("正在执行恢复脚本...")
	restoreScript := filepath.Join(tmpDir, "restore.ps1")
	cmd := exec.Command("powershell", "-File", restoreScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行恢复脚本失败: %v", err)
	}

	fmt.Printf("\n✅ Scoop 配置恢复成功！\n")
	return nil
}
