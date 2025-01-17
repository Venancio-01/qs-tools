package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"qs-tools/internal/config"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 创建 SSH 客户端配置
func createSSHConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: config.DefaultServerUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.DefaultSSHPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// 连接到 SFTP 服务器
func connectSFTP() (*sftp.Client, *ssh.Client, error) {
	sshConfig := createSSHConfig()

	// 连接到 SSH 服务器
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", config.DefaultServerIP), sshConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("连接 SSH 服务器失败: %v", err)
	}

	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, nil, fmt.Errorf("创建 SFTP 客户端失败: %v", err)
	}

	return sftpClient, sshClient, nil
}

// DownloadFromRemote 从远程服务器下载文件
func DownloadFromRemote(component, localFile string) error {
	// 连接到 SFTP 服务器
	sftpClient, sshClient, err := connectSFTP()
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 构建远程文件路径（直接从上传目录下获取）
	remoteFile := filepath.Join(config.DefaultServerPath, fmt.Sprintf("%s_backup", component))
	if runtime.GOOS == "windows" {
		remoteFile += ".zip"
	} else {
		remoteFile += ".tar.gz"
	}
	remoteFile = filepath.ToSlash(remoteFile)

	fmt.Printf("正在从 %s 下载文件...\n", remoteFile)

	// 打开远程文件
	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %v", err)
	}
	defer srcFile.Close()

	// 创建本地文件
	dstFile, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}

	return nil
}

// 检查并创建远程目录
func ensureRemoteDir(client *sftp.Client, path string) error {
	// 先检查父目录
	parentDir := filepath.Dir(path)
	parentInfo, err := client.Stat(parentDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("父目录不存在: %s", parentDir)
		}
		return fmt.Errorf("检查父目录失败: %v", err)
	}
	if !parentInfo.IsDir() {
		return fmt.Errorf("父路径不是目录: %s", parentDir)
	}

	// 检查目标路径
	info, err := client.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 目录不存在，尝试创建
			if err := client.MkdirAll(path); err != nil {
				return fmt.Errorf("创建目录失败: %v", err)
			}
			return nil
		}
		return fmt.Errorf("检查目录失败: %v", err)
	}

	// 路径存在，确保是目录
	if !info.IsDir() {
		return fmt.Errorf("目标路径已存在但不是目录: %s", path)
	}

	return nil
}

// UploadToRemote 上传文件到远程服务器
func UploadToRemote(component, localFile string) error {
	// 连接到 SFTP 服务器
	sftpClient, sshClient, err := connectSFTP()
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 检查上传目录是否存在
	info, err := sftpClient.Stat(config.DefaultServerPath)
	if err != nil {
		return fmt.Errorf("检查上传目录失败: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("上传路径不是目录: %s", config.DefaultServerPath)
	}

	// 构建远程文件路径（直接放在上传目录下）
	remoteFile := filepath.Join(config.DefaultServerPath, fmt.Sprintf("%s_backup", component))
	if runtime.GOOS == "windows" {
		remoteFile += ".zip"
	} else {
		remoteFile += ".tar.gz"
	}
	remoteFile = filepath.ToSlash(remoteFile)

	fmt.Printf("正在上传到 %s...\n", remoteFile)

	// 打开本地文件
	srcFile, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer srcFile.Close()

	// 创建远程文件
	dstFile, err := sftpClient.Create(remoteFile)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %v", err)
	}
	defer dstFile.Close()

	// 复制文件内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	return nil
}

// getRemoteFileName 根据本地文件名获取远程文件名
func getRemoteFileName(localFile string) string {
	base := strings.TrimSuffix(localFile, ".tar.gz")
	base = strings.TrimSuffix(base, ".zip")
	if runtime.GOOS == "windows" {
		return base + ".zip"
	}
	return base + ".tar.gz"
}
