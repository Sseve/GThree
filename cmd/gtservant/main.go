package main

import (
	"GThree/pkg/grpc/gtservant"
	"GThree/pkg/utils"

	"github.com/spf13/viper"
)

func init() {
	utils.InitConfig("gtservant")
	utils.InitRedis()
}

func main() {

	if viper.GetBool("app_daemon") {
		// 以守护进程方式启动
		utils.Daemon()
	}
	// 运行gtservant app
	gtservant.Start()
}
