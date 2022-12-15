package utils

import (
	"context"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	// "github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Db  *mongo.Database
	RDS *redis.Client
)

// 初始化配置
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

// 初始化mongo数据库
func InitDatabase() {
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

	// 创建user集合
	uCollection := Db.Collection("user")
	// 创建zone集合
	// zCollection := Db.Collection("zone")

	// 设置唯一索引
	opts := options.Index().SetUnique(true)
	if _, err = uCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: opts,
	}); err != nil {
		log.Println("设置user集合唯一索引失败", err)
	}

	// if _, err = zCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
	// 	Keys:    bson.D{{Key: "zid", Value: 1}},
	// 	Options: opts,
	// }); err != nil {
	// 	log.Println("设置zone集合唯一索引失败", err)
	// }
}

// 初始化redis
func InitRedis() {
	rds := redis.NewClient(&redis.Options{
		Addr: viper.GetString("rds_url"),
		// Username: viper.GetString("rds_user"),
		// Password: viper.GetString("rds_pass"),
		DB:       viper.GetInt("rds_db"),
		PoolSize: 1000,
	})

	RDS = rds
}
