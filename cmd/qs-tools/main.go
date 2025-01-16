package main

import (
	"qs-tools/internal/cmd"

	"github.com/sirupsen/logrus"
)

func init() {
	// 配置日志
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	cmd.Execute()
}
