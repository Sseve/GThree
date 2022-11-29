package utils

import (
	"context"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Database

func InitConfig(name string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println("get current path failed: ", err)
	}
	if name == "gtmaster" {
		viper.SetConfigName("gtmaster")
	}
	if name == "gtservant" {
		viper.SetConfigName("gtservant")
	}
	viper.AddConfigPath(pwd)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("init config failed: ", err)
	}
	// 修改配置后,无需重启应用
	go func() {
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Printf("配置文件已经更改: %s\n", in.String())

		})
	}()
}

func InitDatabase() {
	// clientOptions := options.Client().ApplyURI(viper.GetString("db_url"))
	// client, err := mongo.Connect(context.TODO(), clientOptions)
	// if err != nil {
	// 	log.Println("connect mongdb failed: ", err)
	// }

	// err = client.Ping(context.TODO(), nil)
	// if err != nil {
	// 	log.Println("check connect mongdb failed: ", err)
	// } else {
	// 	DbClient = client
	// }

	// 连接池
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("db_contect_timeout"))
	defer cancel()

	option := options.Client().ApplyURI(viper.GetString("db_url"))
	option.SetMaxPoolSize(viper.GetUint64("db_pool_size"))
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		log.Println("connect mongdb failed: ", err)
	}
	Db = client.Database(viper.GetString("db_name"))
}
