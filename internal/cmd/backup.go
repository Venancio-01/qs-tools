package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// 默认的服务器配置
const (
	defaultServerIP   = "107.173.165.209"
	defaultServerUser = "root"
	defaultServerPath = "/root/upload/"
	defaultSSHPass    = "#6aL*k5d&2Lg*V"
	defaultHttpPort   = "22814"
)

var backupCmd = &cobra.Command{
	Use:   "backup [component]",
	Short: "备份配置文件",
	Long: `备份各种工具和软件的配置文件。
目前支持的组件：
  - fish: 备份 Fish Shell 的配置文件
  - scoop: 备份 Scoop 包管理器的配置（Windows）

备份文件将自动上传到配置的远程服务器。`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请指定要备份的组件")
			return
		}

		switch args[0] {
		case "fish":
			backupFish()
		case "scoop":
			backupScoop()
		default:
			fmt.Printf("不支持的组件: %s\n", args[0])
		}
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply [component]",
	Short: "应用配置文件",
	Long: `从远程服务器下载并应用配置文件。
目前支持的组件：
  - fish: 应用 Fish Shell 的配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请指定要应用的组件")
			return
		}

		switch args[0] {
		case "fish":
			applyFish()
		default:
			fmt.Printf("不支持的组件: %s\n", args[0])
		}
	},
}

func init() {
	RootCmd.AddCommand(backupCmd)
	RootCmd.AddCommand(applyCmd)
}

func uploadToRemoteServer(localFile string) error {
	// 使用 sshpass 和 scp 上传文件
	fmt.Println("\n开始上传备份文件到远程服务器...")

	// 检查 sshpass 是否安装
	if _, err := exec.LookPath("sshpass"); err != nil {
		fmt.Println("正在安装 sshpass...")
		installCmd := exec.Command("sudo", "apt-get", "install", "-y", "sshpass")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("安装 sshpass 失败: %v", err)
		}
	}

	// 构建远程路径
	remotePath := fmt.Sprintf("%s@%s:%s", defaultServerUser, defaultServerIP, defaultServerPath)

	// 使用 sshpass 执行 scp 命令
	cmd := exec.Command("sshpass", "-p", defaultSSHPass, "scp", "-o", "StrictHostKeyChecking=no", localFile, remotePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	fmt.Printf("✅ 备份文件已成功上传到：%s:%s\n", defaultServerIP, defaultServerPath)
	return nil
}

func backupFish() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return
	}

	fishConfigDir := filepath.Join(homeDir, ".config", "fish")
	if _, err := os.Stat(fishConfigDir); os.IsNotExist(err) {
		fmt.Printf("Fish 配置目录不存在: %s\n", fishConfigDir)
		return
	}

	// 创建备份目录
	backupDir := filepath.Join(homeDir, "qs-tools-backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fmt.Printf("创建备份目录失败: %v\n", err)
		return
	}

	// 创建备份文件名
	backupFile := filepath.Join(backupDir, "fish-config.tar.gz")

	// 创建压缩文件
	file, err := os.Create(backupFile)
	if err != nil {
		fmt.Printf("创建备份文件失败: %v\n", err)
		return
	}

	// 创建 gzip writer
	gw := gzip.NewWriter(file)
	// 创建 tar writer
	tw := tar.NewWriter(gw)

	fileCount := 0
	var totalSize int64

	// 遍历 Fish 配置目录
	err = filepath.Walk(fishConfigDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历目录失败: %v", err)
		}

		// 获取相对路径
		relPath, err := filepath.Rel(fishConfigDir, path)
		if err != nil {
			return fmt.Errorf("获取相对路径失败: %v", err)
		}

		// 创建 tar header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return fmt.Errorf("创建文件头失败: %v", err)
		}
		header.Name = filepath.Join("fish", relPath)

		// 写入 header
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("写入文件头失败: %v", err)
		}

		// 如果是普通文件，写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("打开文件失败 %s: %v", path, err)
			}

			if _, err := io.Copy(tw, file); err != nil {
				file.Close()
				return fmt.Errorf("写入文件内容失败 %s: %v", path, err)
			}
			file.Close()

			totalSize += info.Size() // 使用文件的实际大小
			fileCount++
			fmt.Printf("已备份: %s (%.2f KB)\n", relPath, float64(info.Size())/1024)
		}

		return nil
	})

	if err != nil {
		tw.Close()
		gw.Close()
		file.Close()
		os.Remove(backupFile) // 清理损坏的文件
		fmt.Printf("备份过程中出错: %v\n", err)
		return
	}

	// 按顺序关闭写入器
	if err := tw.Close(); err != nil {
		gw.Close()
		file.Close()
		os.Remove(backupFile)
		fmt.Printf("关闭 tar writer 失败: %v\n", err)
		return
	}

	if err := gw.Close(); err != nil {
		file.Close()
		os.Remove(backupFile)
		fmt.Printf("关闭 gzip writer 失败: %v\n", err)
		return
	}

	if err := file.Close(); err != nil {
		os.Remove(backupFile)
		fmt.Printf("关闭文件失败: %v\n", err)
		return
	}

	// 验证备份文件
	stat, err := os.Stat(backupFile)
	if err != nil {
		fmt.Printf("获取备份文件信息失败: %v\n", err)
		return
	}

	if stat.Size() == 0 {
		fmt.Println("错误：备份文件大小为 0")
		os.Remove(backupFile)
		return
	}

	// 验证备份文件是否可读
	testFile, err := os.Open(backupFile)
	if err != nil {
		fmt.Printf("无法打开备份文件进行验证: %v\n", err)
		os.Remove(backupFile)
		return
	}
	defer testFile.Close()

	// 尝试读取 gzip 头
	testGz, err := gzip.NewReader(testFile)
	if err != nil {
		fmt.Printf("备份文件格式无效: %v\n", err)
		os.Remove(backupFile)
		return
	}
	testGz.Close()

	fmt.Printf("\n✅ Fish 配置文件备份成功！\n")
	fmt.Printf("备份文件位置: %s\n", backupFile)
	fmt.Printf("备份文件大小: %.2f KB\n", float64(stat.Size())/1024)
	fmt.Printf("包含文件数量: %d\n", fileCount)
	fmt.Printf("原始文件总大小: %.2f KB\n", float64(totalSize)/1024)

	// 显示备份内容概要
	fmt.Println("\n备份内容包括：")
	fmt.Printf("- 配置目录: %s\n", fishConfigDir)

	// 列出主要配置文件
	mainConfigs := []string{"config.fish", "functions", "completions"}
	for _, conf := range mainConfigs {
		path := filepath.Join(fishConfigDir, conf)
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("- %s\n", conf)
		}
	}

	// 上传到远程服务器
	if err := uploadToRemoteServer(backupFile); err != nil {
		fmt.Printf("\n⚠️ %v\n", err)
		os.Remove(backupFile)
		return
	}

	// 清理本地备份文件
	if err := os.Remove(backupFile); err != nil {
		fmt.Printf("\n清理本地备份文件失败: %v\n", err)
	}
}

func applyFish() {
	fmt.Println("开始应用 Fish 配置...")

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "fish-config")
	if err != nil {
		fmt.Printf("创建临时目录失败: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	// 下载配置文件
	configUrl := fmt.Sprintf("http://%s:%s/fish-config.tar.gz", defaultServerIP, defaultHttpPort)
	fmt.Printf("正在从 %s 下载配置文件...\n", configUrl)

	resp, err := http.Get(configUrl)
	if err != nil {
		fmt.Printf("下载配置文件失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("下载失败，服务器返回状态码: %d\n", resp.StatusCode)
		return
	}

	// 保存下载的文件
	tmpFile := filepath.Join(tmpDir, "config.tar.gz")
	out, err := os.Create(tmpFile)
	if err != nil {
		fmt.Printf("创建临时文件失败: %v\n", err)
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		fmt.Printf("保存配置文件失败: %v\n", err)
		return
	}

	// 解压配置文件
	fmt.Println("正在解压配置文件...")

	// 打开下载的压缩文件
	f, err := os.Open(tmpFile)
	if err != nil {
		fmt.Printf("打开配置文件失败: %v\n", err)
		return
	}
	defer f.Close()

	// 创建 gzip reader
	gr, err := gzip.NewReader(f)
	if err != nil {
		fmt.Printf("读取压缩文件失败: %v\n", err)
		return
	}
	defer gr.Close()

	// 创建 tar reader
	tr := tar.NewReader(gr)

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return
	}

	// 备份现有配置
	fishConfigDir := filepath.Join(homeDir, ".config", "fish")
	if _, err := os.Stat(fishConfigDir); err == nil {
		backupDir := filepath.Join(homeDir, ".config", "fish.bak")
		if err := os.RemoveAll(backupDir); err != nil {
			fmt.Printf("删除旧的备份失败: %v\n", err)
			return
		}
		if err := os.Rename(fishConfigDir, backupDir); err != nil {
			fmt.Printf("备份现有配置失败: %v\n", err)
			return
		}
		fmt.Printf("已备份现有配置到: %s\n", backupDir)
	}

	// 创建新的配置目录
	if err := os.MkdirAll(fishConfigDir, 0755); err != nil {
		fmt.Printf("创建配置目录失败: %v\n", err)
		return
	}

	// 解压文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("读取压缩文件内容失败: %v\n", err)
			return
		}

		// 跳过 "fish/" 前缀
		name := header.Name
		if strings.HasPrefix(name, "fish/") {
			name = name[5:]
		}
		if name == "" {
			continue
		}

		target := filepath.Join(fishConfigDir, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				return
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("创建目录失败: %v\n", err)
				return
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				fmt.Printf("创建文件失败: %v\n", err)
				return
			}
			defer f.Close()

			if _, err := io.Copy(f, tr); err != nil {
				fmt.Printf("写入文件失败: %v\n", err)
				return
			}
		}
	}

	fmt.Println("\n✅ Fish 配置应用成功！")
	fmt.Println("\n提示：")
	fmt.Println("1. 重新启动 Fish Shell 以应用新的配置")
	fmt.Println("2. 如果出现问题，可以从 ~/.config/fish.bak 恢复原有配置")
}

func backupScoop() error {
	fmt.Println("开始备份 Scoop 配置...")

	if runtime.GOOS != "windows" {
		return fmt.Errorf("Scoop 仅支持 Windows 系统")
	}

	// 1. 获取 Scoop 目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户主目录失败: %v", err)
	}
	scoopDir := filepath.Join(homeDir, "scoop")

	// 检查 Scoop 是否已安装
	if _, err := os.Stat(scoopDir); os.IsNotExist(err) {
		return fmt.Errorf("未找到 Scoop 安装目录")
	}

	// 2. 创建临时目录
	tmpDir, err := os.MkdirTemp("", "scoop-backup")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 3. 导出已安装的应用列表
	fmt.Println("\n1. 导出已安装的应用列表...")
	appsFile := filepath.Join(tmpDir, "apps.txt")
	listCmd := exec.Command("scoop", "list")
	appsOutput, err := listCmd.Output()
	if err != nil {
		return fmt.Errorf("导出应用列表失败: %v", err)
	}
	if err := os.WriteFile(appsFile, appsOutput, 0644); err != nil {
		return fmt.Errorf("保存应用列表失败: %v", err)
	}

	// 4. 导出软件源列表
	fmt.Println("2. 导出软件源列表...")
	bucketsFile := filepath.Join(tmpDir, "buckets.txt")
	bucketsCmd := exec.Command("scoop", "bucket", "list")
	bucketsOutput, err := bucketsCmd.Output()
	if err != nil {
		return fmt.Errorf("导出软件源列表失败: %v", err)
	}
	if err := os.WriteFile(bucketsFile, bucketsOutput, 0644); err != nil {
		return fmt.Errorf("保存软件源列表失败: %v", err)
	}

	// 5. 备份配置文件
	fmt.Println("3. 备份配置文件...")
	configDir := filepath.Join(scoopDir, "config")
	if _, err := os.Stat(configDir); err == nil {
		configBackupDir := filepath.Join(tmpDir, "config")
		if err := copyDir(configDir, configBackupDir); err != nil {
			return fmt.Errorf("备份配置文件失败: %v", err)
		}
	}

	// 6. 生成恢复脚本
	fmt.Println("4. 生成恢复脚本...")
	restoreScript := `@echo off
echo 开始恢复 Scoop 配置...

echo 1. 安装 Scoop...
powershell -Command "Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force; iwr -useb get.scoop.sh | iex"

echo 2. 添加软件源...
for /F "tokens=*" %%b in (buckets.txt) do (
    scoop bucket add %%b
)

echo 3. 恢复配置文件...
if exist config (
    xcopy /E /I /Y config %USERPROFILE%\scoop\config
)

echo 4. 安装应用...
for /F "tokens=*" %%a in (apps.txt) do (
    scoop install %%a
)

echo ✅ Scoop 配置恢复完成！
pause
`
	restoreFile := filepath.Join(tmpDir, "restore.bat")
	if err := os.WriteFile(restoreFile, []byte(restoreScript), 0644); err != nil {
		return fmt.Errorf("生成恢复脚本失败: %v", err)
	}

	// 7. 创建备份文件
	fmt.Println("5. 创建备份文件...")
	backupPath := filepath.Join(tmpDir, "scoop-backups.zip")

	// 使用 PowerShell 的 Compress-Archive 命令创建 ZIP 文件
	zipCmd := exec.Command("powershell", "-Command",
		fmt.Sprintf("Compress-Archive -Path '%s\\*' -DestinationPath '%s'", tmpDir, backupPath))
	zipCmd.Stdout = os.Stdout
	zipCmd.Stderr = os.Stderr
	if err := zipCmd.Run(); err != nil {
		return fmt.Errorf("创建备份文件失败: %v", err)
	}

	// 8. 上传到远程服务器
	fmt.Println("\n6. 上传备份文件到远程服务器...")

	// 使用 SFTP 上传文件
	sftpCmd := exec.Command("sftp", "-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("%s@%s", defaultServerUser, defaultServerIP))

	// 构建 SFTP 命令
	sftpCommands := fmt.Sprintf("cd %s\nput %s\nbye\n", defaultServerPath, backupPath)
	sftpCmd.Stdin = strings.NewReader(sftpCommands)
	sftpCmd.Stdout = os.Stdout
	sftpCmd.Stderr = os.Stderr

	if err := sftpCmd.Run(); err != nil {
		return fmt.Errorf("上传文件失败: %v", err)
	}

	fmt.Printf("\n✅ Scoop 配置备份完成！\n")
	fmt.Printf("备份文件已上传到：%s:%s\n", defaultServerIP, defaultServerPath)
	return nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
