package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

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
	tmpDir, err := os.MkdirTemp("", "fish-backup")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建压缩文件
	backupFile := filepath.Join(tmpDir, "fish_backup.tar.gz")
	fmt.Println("正在压缩配置文件...")
	cmd := exec.Command("tar", "czf", backupFile, "-C", filepath.Dir(configDir), filepath.Base(configDir))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("压缩配置文件失败: %v", err)
	}

	// 上传到远程服务器
	fmt.Println("正在上传到远程服务器...")
	remotePath := fmt.Sprintf("%s@%s:%s/fish", defaultServerUser, defaultServerIP, defaultServerPath)
	cmd = exec.Command("sshpass", "-p", defaultSSHPass, "scp", "-o", "StrictHostKeyChecking=no",
		backupFile, remotePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("上传配置文件失败: %v", err)
	}

	fmt.Printf("\n✅ Fish Shell 配置备份成功！\n")
	fmt.Printf("备份文件已上传到: %s\n", remotePath)
	return nil
} 
