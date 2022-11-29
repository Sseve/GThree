package dto

import (
	"GThree/pkg/models"
	"GThree/pkg/utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "gopkg.in/mgo.v2"
)

type DUser struct {
	Name       string
	Password   string
	Desc       string
	Roles      []string
	Avatar     string
	CreateTime string
	UpdateTime string
}

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

func UptUserToDb() {

}

func SelectUserFromDb(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := utils.Db.Collection("user").FindOne(ctx, nil)
	fmt.Println("select from db: ", result)
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
