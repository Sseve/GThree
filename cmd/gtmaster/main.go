package main

import (
	"GThree/pkg/route"
	"GThree/pkg/utils"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func init() {
	utils.InitConfig("gtmaster")
	utils.InitDatabase()
	utils.InitRedis()
}

func main() {
	if viper.GetBool("app_daemon") {
		// 以守护进程方式启动
		utils.Daemon()
	}

	// 运行gtmaster app
	startApp()
}

// 运行gin实例
func startApp() {
	serve := http.Server{
		Addr:         viper.GetString("app_addr"),
		Handler:      route.GetRouter(),
		ReadTimeout:  viper.GetDuration("app_read_timeout") * time.Second,
		WriteTimeout: viper.GetDuration("app_write_timeout") * time.Second,
	}

	go func() {
		log.Println("server start on: ", viper.GetString("app_addr"))
		if err := serve.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server start failed: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := serve.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed: ", err)
	} else {
		log.Println("server exit success")
	}
}
