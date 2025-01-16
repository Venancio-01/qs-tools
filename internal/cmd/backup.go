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

func backupScoop() {
	if runtime.GOOS != "windows" {
		fmt.Println("Scoop 备份功能仅支持 Windows 系统")
		return
	}

	// 检查是否安装了 scoop
	if _, err := exec.LookPath("scoop"); err != nil {
		fmt.Println("未检测到 Scoop，请先安装 Scoop")
		return
	}

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v\n", err)
		return
	}

	// 创建备份目录
	backupDir := filepath.Join(homeDir, "qs-tools-backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fmt.Printf("创建备份目录失败: %v\n", err)
		return
	}

	// 创建临时目录存放配置文件
	tmpDir := filepath.Join(backupDir, "scoop-backup")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		fmt.Printf("创建临时目录失败: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("开始备份 Scoop 配置...")

	// 1. 导出已安装的应用列表
	fmt.Println("\n1. 导出已安装的应用列表...")
	installedApps := filepath.Join(tmpDir, "installed-apps.txt")
	listCmd := exec.Command("scoop", "list")
	output, err := listCmd.Output()
	if err != nil {
		fmt.Printf("获取已安装应用列表失败: %v\n", err)
		return
	}
	if err := os.WriteFile(installedApps, output, 0644); err != nil {
		fmt.Printf("保存应用列表失败: %v\n", err)
		return
	}

	// 2. 导出 bucket 列表
	fmt.Println("2. 导出软件源列表...")
	bucketList := filepath.Join(tmpDir, "bucket-list.txt")
	bucketCmd := exec.Command("scoop", "bucket", "list")
	output, err = bucketCmd.Output()
	if err != nil {
		fmt.Printf("获取软件源列表失败: %v\n", err)
		return
	}
	if err := os.WriteFile(bucketList, output, 0644); err != nil {
		fmt.Printf("保存软件源列表失败: %v\n", err)
		return
	}

	// 3. 复制 Scoop 配置文件
	fmt.Println("3. 备份配置文件...")
	scoopConfigDir := filepath.Join(homeDir, "scoop", "config")
	if _, err := os.Stat(scoopConfigDir); err == nil {
		configBackupDir := filepath.Join(tmpDir, "config")
		if err := os.MkdirAll(configBackupDir, 0755); err != nil {
			fmt.Printf("创建配置备份目录失败: %v\n", err)
			return
		}

		// 复制所有配置文件
		files, err := os.ReadDir(scoopConfigDir)
		if err != nil {
			fmt.Printf("读取配置目录失败: %v\n", err)
			return
		}

		for _, file := range files {
			srcPath := filepath.Join(scoopConfigDir, file.Name())
			dstPath := filepath.Join(configBackupDir, file.Name())
			input, err := os.ReadFile(srcPath)
			if err != nil {
				fmt.Printf("读取配置文件 %s 失败: %v\n", file.Name(), err)
				continue
			}
			if err := os.WriteFile(dstPath, input, 0644); err != nil {
				fmt.Printf("保存配置文件 %s 失败: %v\n", file.Name(), err)
				continue
			}
		}
	}

	// 4. 创建恢复脚本
	fmt.Println("4. 生成恢复脚本...")
	restoreScript := filepath.Join(tmpDir, "restore.ps1")
	scriptContent := `# Scoop 配置恢复脚本
Write-Host "开始恢复 Scoop 配置..."

# 1. 检查并安装 Scoop
if (!(Get-Command scoop -ErrorAction SilentlyContinue)) {
    Write-Host "未检测到 Scoop，正在安装..."
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
    Invoke-RestScript -Uri get.scoop.sh | Invoke-Expression
}

# 2. 添加软件源
Write-Host "正在添加软件源..."
Get-Content "bucket-list.txt" | ForEach-Object {
    Write-Host "添加软件源: $_"
    scoop bucket add $_
}

# 3. 安装应用
Write-Host "正在安装应用..."
Get-Content "installed-apps.txt" | ForEach-Object {
    if ($_ -match "^(\S+)\s+") {
        $app = $matches[1]
        Write-Host "安装应用: $app"
        scoop install $app
    }
}

# 4. 恢复配置文件
$configDir = Join-Path $env:USERPROFILE "scoop\config"
if (Test-Path "config") {
    Write-Host "正在恢复配置文件..."
    if (!(Test-Path $configDir)) {
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    }
    Copy-Item "config\*" $configDir -Force
}

Write-Host "✅ Scoop 配置恢复完成！"
`
	if err := os.WriteFile(restoreScript, []byte(scriptContent), 0644); err != nil {
		fmt.Printf("生成恢复脚本失败: %v\n", err)
		return
	}

	// 5. 创建压缩文件
	backupFile := filepath.Join(backupDir, "scoop-config.tar.gz")
	fmt.Println("5. 创建备份文件...")

	// 创建 tar.gz 文件
	file, err := os.Create(backupFile)
	if err != nil {
		fmt.Printf("创建备份文件失败: %v\n", err)
		return
	}

	gw := gzip.NewWriter(file)
	tw := tar.NewWriter(gw)

	// 遍历临时目录中的所有文件
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(tmpDir, path)
		if err != nil {
			return err
		}

		// 跳过目录本身
		if relPath == "." {
			return nil
		}

		// 创建 header
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		// 写入 header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// 如果是普通文件，写入内容
		if !info.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if _, err := tw.Write(data); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		tw.Close()
		gw.Close()
		file.Close()
		fmt.Printf("打包备份文件失败: %v\n", err)
		return
	}

	// 关闭所有写入器
	if err := tw.Close(); err != nil {
		gw.Close()
		file.Close()
		fmt.Printf("关闭 tar writer 失败: %v\n", err)
		return
	}
	if err := gw.Close(); err != nil {
		file.Close()
		fmt.Printf("关闭 gzip writer 失败: %v\n", err)
		return
	}
	if err := file.Close(); err != nil {
		fmt.Printf("关闭文件失败: %v\n", err)
		return
	}

	// 6. 上传到远程服务器
	fmt.Println("6. 上传备份文件到远程服务器...")
	if err := uploadToRemoteServer(backupFile); err != nil {
		fmt.Printf("\n⚠️ %v\n", err)
		return
	}

	// 7. 清理本地备份文件
	if err := os.Remove(backupFile); err != nil {
		fmt.Printf("\n清理本地备份文件失败: %v\n", err)
		return
	}

	fmt.Println("\n✅ Scoop 配置备份成功！")
	fmt.Println("\n备份内容包括：")
	fmt.Println("1. 已安装的应用列表")
	fmt.Println("2. 已添加的软件源列表")
	fmt.Println("3. Scoop 配置文件")
	fmt.Println("4. 恢复脚本 (restore.ps1)")
}
