package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

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
	tmpDir, err := os.MkdirTemp("", "nvim-backup")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建压缩文件
	backupFile := filepath.Join(tmpDir, "nvim_backup.tar.gz")

	// 创建压缩文件
	fmt.Println("正在压缩配置文件...")
	if runtime.GOOS == "windows" {
		// Windows 使用 PowerShell 压缩
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Compress-Archive -Path '%s\\*' -DestinationPath '%s'",
				configDir, backupFile[:len(backupFile)-7] + "_backup.zip"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("压缩配置文件失败: %v", err)
		}
		backupFile = backupFile[:len(backupFile)-7] + "_backup.zip"
	} else {
		// Linux 使用 tar
		cmd := exec.Command("tar", "czf", backupFile, "-C", filepath.Dir(configDir), filepath.Base(configDir))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("压缩配置文件失败: %v", err)
		}
	}

	// 上传到远程服务器
	fmt.Println("正在上传到远程服务器...")
	remotePath := fmt.Sprintf("%s@%s:%s/nvim", defaultServerUser, defaultServerIP, defaultServerPath)

	if runtime.GOOS == "windows" {
		// Windows 使用 SFTP 批处理命令
		sftpCommands := fmt.Sprintf("cd %s\nput %s\nbye\n", defaultServerPath+"/nvim", backupFile)
		cmd := exec.Command("sftp", "-b", "-", fmt.Sprintf("%s@%s", defaultServerUser, defaultServerIP))
		cmd.Stdin = strings.NewReader(sftpCommands)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("上传配置文件失败: %v", err)
		}
	} else {
		// Linux 使用 sshpass + scp
		cmd := exec.Command("sshpass", "-p", defaultSSHPass, "scp", "-o", "StrictHostKeyChecking=no",
			backupFile, remotePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("上传配置文件失败: %v", err)
		}
	}

	fmt.Printf("\n✅ Neovim 配置备份成功！\n")
	fmt.Printf("备份文件已上传到: %s\n", remotePath)
	return nil
}
