package utils

import (
	"fmt"
	"os"
)

// CreateTempDir 创建临时目录并返回清理函数
func CreateTempDir(prefix string) (dir string, cleanup func(), err error) {
	tmpDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", nil, fmt.Errorf("创建临时目录失败: %v", err)
	}

	cleanup = func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup, nil
} 
