package dto

import (
	"GThree/pkg/models"
	"GThree/pkg/utils"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

///////////////////////// mongo /////////////////////////
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

//////////////////////////// redis /////////////////////////////
// 操作入库
func SetZResToRds(name string, resp interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := utils.RDS.Set(ctx, name, resp, -1).Result()
	if err != nil {
		utils.Logger.Error("Rds记录结果数据出错: ", err)
	}
	log.Println("Rds记录结果数据成功: ", res)
}

// 操作出库
func GetZResToRds(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := utils.RDS.Get(ctx, name).Result()
	if err != nil {
		utils.Logger.Info("Rds记录结果数据出错: ", err)
		return "", nil
	}
	if res != "" {
		utils.RDS.Del(ctx, name)
		utils.Logger.Info("Rds记录结果数据删除: ", name)
	}
	return res, nil
}
