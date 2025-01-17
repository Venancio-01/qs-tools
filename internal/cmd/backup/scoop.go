package backup

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
	Short: "备份 Scoop 配置",
	Long:  `备份 Scoop 配置文件并上传到远程服务器。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return backupScoop()
	},
}

func init() {
	BackupCmd.AddCommand(scoopCmd)
}

func backupScoop() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Scoop 仅支持 Windows 系统")
	}

	fmt.Println("开始备份 Scoop 配置...")

	// 创建临时目录
	tmpDir, cleanup, err := utils.CreateTempDir("scoop-backup")
	if err != nil {
		return err
	}
	defer cleanup()

	// 导出已安装的应用列表
	fmt.Println("导出已安装的应用列表...")
	appsFile := filepath.Join(tmpDir, "apps.txt")
	cmd := exec.Command("scoop", "list")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("获取应用列表失败: %v", err)
	}
	if err := os.WriteFile(appsFile, output, 0644); err != nil {
		return fmt.Errorf("写入应用列表失败: %v", err)
	}

	// 导出软件源列表
	fmt.Println("导出软件源列表...")
	bucketsFile := filepath.Join(tmpDir, "buckets.txt")
	cmd = exec.Command("scoop", "bucket", "list")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("获取软件源列表失败: %v", err)
	}
	if err := os.WriteFile(bucketsFile, output, 0644); err != nil {
		return fmt.Errorf("写入软件源列表失败: %v", err)
	}

	// 生成恢复脚本
	fmt.Println("生成恢复脚本...")
	restoreScript := filepath.Join(tmpDir, "restore.ps1")
	scriptContent := `# 添加软件源
Get-Content buckets.txt | ForEach-Object {
    scoop bucket add $_
}

# 安装应用
Get-Content apps.txt | ForEach-Object {
    scoop install $_
}`
	if err := os.WriteFile(restoreScript, []byte(scriptContent), 0644); err != nil {
		return fmt.Errorf("生成恢复脚本失败: %v", err)
	}

	// 创建压缩文件
	backupFile := filepath.Join(tmpDir, "scoop_backup.zip")
	if err := utils.CompressDir(tmpDir, backupFile); err != nil {
		return err
	}

	// 上传到远程服务器
	if err := utils.UploadToRemote("scoop", backupFile); err != nil {
		return err
	}

	fmt.Printf("\n✅ Scoop 配置备份成功！\n")
	return nil
}
