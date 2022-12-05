package utils

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	if _, err := os.Stat("logs"); err != nil {
		if !os.IsExist(err) {
			if err := os.Mkdir("logs", 0755); err != nil {
				log.Println("创建日志文件夹成功")
			} else {
				log.Println("创建日志文件夹失败")
			}
		}
	}
	fd, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("创建日志文件失败", err)
	}
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(fd)
}
