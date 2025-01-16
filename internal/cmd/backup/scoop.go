package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func backupScoop() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("Scoop 仅支持 Windows 系统")
	}

	fmt.Println("开始备份 Scoop 配置...")

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "scoop-backup")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

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
	fmt.Println("正在压缩配置文件...")
	cmd = exec.Command("powershell", "-Command",
		fmt.Sprintf("Compress-Archive -Path '%s\\*' -DestinationPath '%s'",
			tmpDir, backupFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("压缩配置文件失败: %v", err)
	}

	// 上传到远程服务器
	fmt.Println("正在上传到远程服务器...")
	remotePath := fmt.Sprintf("%s@%s:%s/scoop", defaultServerUser, defaultServerIP, defaultServerPath)

	// Windows 使用 SFTP 批处理命令
	sftpCommands := fmt.Sprintf("cd %s\nput %s\nbye\n", defaultServerPath+"/scoop", backupFile)
	cmd = exec.Command("sftp", "-b", "-", fmt.Sprintf("%s@%s", defaultServerUser, defaultServerIP))
	cmd.Stdin = strings.NewReader(sftpCommands)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("上传配置文件失败: %v", err)
	}

	fmt.Printf("\n✅ Scoop 配置备份成功！\n")
	fmt.Printf("备份文件已上传到: %s\n", remotePath)
	return nil
} 
