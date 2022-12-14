package main

import (
	"GThree/pkg/grpc/gtservant"
	"GThree/pkg/utils"
)

func init() {
	utils.InitConfig("gtservant")
	utils.InitRedis()
}

func main() {
	gtservant.Start()
}
