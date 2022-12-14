package dto

import (
	"GThree/pkg/models"
	"GThree/pkg/utils"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// 区服数据库模型
type DZone struct {
	Zid        string
	Ip         string
	Name       string
	Closed     bool // 是否关服
	CreateTime string
	UpdateTime string
}

// 添加区服信息
func AddZoneToDb(zones models.ZoneOpt) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	documents := make([]interface{}, 0, len(zones.Zone))
	for _, zone := range zones.Zone {
		documents = append(documents, DZone{
			Zid:        zone.Zid,
			Name:       zone.Name,
			Ip:         zone.Ip,
			Closed:     zone.Closed,
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
			UpdateTime: "",
		})
	}
	if _, err := utils.Db.Collection("zone").InsertMany(ctx, documents); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// 删除区服信息
func DelZoneFromDb(zid, name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fiter := bson.M{"zid": zid, "name": name}
	result, err := utils.Db.Collection("zone").DeleteOne(ctx, fiter)
	if err != nil {
		return false
	}
	if result.DeletedCount == 0 {
		return false
	}
	return true
}

// 更新区服信息
func UptZoneToDb() {

}

// 查询区服信息
func SelectZoneFromDb(zid, name string) (*DZone, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fiter := bson.M{"zid": zid, "name": name}
	// opt := options.FindOne().SetProjection(bson.M{"name": 1})
	var zone DZone
	err := utils.Db.Collection("user").FindOne(ctx, fiter).Decode(&zone)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}
