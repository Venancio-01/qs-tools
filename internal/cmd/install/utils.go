package install

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// isDebianBased 检查是否为基于 Debian 的系统
func isDebianBased() bool {
	// 检查 /etc/os-release 文件
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}

	content := string(data)
	return strings.Contains(content, "Ubuntu") ||
		strings.Contains(content, "Debian") ||
		strings.Contains(content, "Kylin") ||
		strings.Contains(content, "LinuxMint") ||
		strings.Contains(content, "Pop!_OS")
}

// isKylin 检查是否为 Kylin 系统
func isKylin() bool {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return false
	}

	return strings.Contains(string(data), "Kylin")
}

// downloadFile 下载文件到指定路径
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractTarGz 解压 tar.gz 文件
func extractTarGz(archivePath, destPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
		}
	}

	return nil
} 
