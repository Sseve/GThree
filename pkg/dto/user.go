package dto

import (
	"GThree/pkg/models"
	"GThree/pkg/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "gopkg.in/mgo.v2"
)

// 用户入库数据库模型
type DUser struct {
	Name       string
	Password   string
	Desc       string
	Roles      []string
	Avatar     string
	CreateTime string
	UpdateTime string
}

// 用户出库数据模型
type DOUser struct {
	Name       string
	Desc       string
	Roles      []string
	Avatar     string
	CreateTime string
	UpdateTime string
}

// 添加用户
func AddUserToDb(user models.UserAdd) bool {
	u := DUser{
		Name:       user.Name,
		Password:   utils.HashAndSalt([]byte(user.PassOne)),
		Desc:       user.Desc,
		Roles:      user.Roles,
		Avatar:     user.Avatar,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: "",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := utils.Db.Collection("user").InsertOne(ctx, u); err != nil {
		return false
	}
	return true
}

// 删除用户
func DelUserFromDb(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fiter := bson.M{"name": name}
	result, err := utils.Db.Collection("user").DeleteOne(ctx, fiter)
	if err != nil {
		return false
	}
	if result.DeletedCount == 0 {
		return false
	}
	return true
}

// 修改用户数据
func UptUserToDb(name string, user models.UserMod) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"name": name}
	umap := make(bson.M)
	if user.Name != "" {
		umap["name"] = user.Name
	}
	if user.Password != "" {
		umap["password"] = utils.HashAndSalt([]byte(user.Password))
	}
	if user.Desc != "" {
		umap["desc"] = user.Desc
	}
	if user.Roles != nil {
		umap["roles"] = user.Roles
	}
	if user.Avatar != "" {
		umap["avatar"] = user.Avatar
	}
	umap["updatetime"] = time.Now().Format("2006-01-02 15:04:05")
	upt := bson.M{"$set": umap}
	var u DUser
	if err := utils.Db.Collection("user").FindOneAndUpdate(ctx, filter, upt).Decode(&u); err != nil {
		return false
	}
	return true
}

// 查询多个用户
func SelectUserFromDb() ([]*DOUser, error) {
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := utils.Db.Collection("user").Find(ctx1, bson.D{})
	if err != nil {
		return nil, err
	}
	var user []*DOUser
	defer cur.Close(ctx1)
	ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = cur.All(ctx2, &user); err != nil {
		return nil, err
	}
	return user, nil
}

func CheckUserFromDb(name, password string) bool {
	var u DUser
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fiter := bson.M{"name": name}
	err := utils.Db.Collection("user").FindOne(ctx, fiter).Decode(&u)
	if err != nil {
		return false
	}
	if !utils.ComparePassword(u.Password, password) {
		return false
	}
	return true
}

// 查询单个用户
func GetUserFromDb(name string) (*DOUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fiter := bson.M{"name": name}
	// opt := options.FindOne().SetProjection(bson.M{"name": 1})
	var user DOUser
	err := utils.Db.Collection("user").FindOne(ctx, fiter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
