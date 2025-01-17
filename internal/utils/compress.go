package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// CompressDir 压缩目录
func CompressDir(sourceDir, targetFile string) error {
	fmt.Println("正在压缩文件...")

	if runtime.GOOS == "windows" {
		// Windows 使用 PowerShell 压缩
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Compress-Archive -Path '%s\\*' -DestinationPath '%s' -Force",
				sourceDir, targetFile))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("压缩文件失败: %v", err)
		}
	} else {
		// Linux 使用 tar
		cmd := exec.Command("tar", "czf", targetFile,
			"-C", filepath.Dir(sourceDir), filepath.Base(sourceDir))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("压缩文件失败: %v", err)
		}
	}

	return nil
}

// ExtractFile 解压文件
func ExtractFile(sourceFile, targetDir string) error {
	fmt.Println("正在解压文件...")

	if runtime.GOOS == "windows" {
		// Windows 使用 PowerShell 解压
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Expand-Archive -Path '%s' -DestinationPath '%s' -Force",
				sourceFile, targetDir))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("解压文件失败: %v", err)
		}
	} else {
		// Linux 使用 tar
		cmd := exec.Command("tar", "xzf", sourceFile, "-C", targetDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("解压文件失败: %v", err)
		}
	}

	return nil
} 
