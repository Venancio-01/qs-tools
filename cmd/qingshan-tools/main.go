package main

import (
	"log"

	"qingshan-tools/internal/cmd"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// 配置日志
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	cmd.Execute()
}
