package api

import (
	"GThree/pkg/dto"
	"GThree/pkg/models"
	"GThree/pkg/utils"

	"github.com/gin-gonic/gin"
)

type user struct {
	MUSign models.UserSign
	MUAdd  models.UserAdd
	MUUpt  models.UserUpt
}

func Newuser() *user {
	return new(user)
}

func (u *user) Sign(ctx *gin.Context) {
	// 获取接口数据
	if err := ctx.BindJSON(&u.MUSign); err != nil {
		utils.RespFalured(ctx, "获取用户登录api接口数据失败", nil)
		return
	}
	// 校验用户数据
	if !dto.CheckUserFromDb(u.MUSign.Name, u.MUSign.Password) {
		utils.RespFalured(ctx, "用户名或密码错误", nil)
		return
	}
	utils.Logger.Info("用户: ", u.MUSign.Name, "登录成功")
	// 创建token
	token, err := utils.CreateToken(u.MUSign.Name)
	if err != nil {
		utils.RespFalured(ctx, "创建token失败", nil)
		return
	}
	// 返回结果
	utils.RespSuccess(ctx, "登录成功", token)
}

// 添加用户
func (u *user) Add(ctx *gin.Context) {
	if err := ctx.BindJSON(&u.MUAdd); err != nil {
		utils.RespFalured(ctx, "获取添加用户api接口数据失败", nil)
		return
	}
	if !dto.AddUserToDb(u.MUAdd) {
		utils.RespFalured(ctx, "添加用户失败,或许用户已经存在", nil)
		return
	}
	utils.Logger.Info("用户: ", u.MUAdd.Name, "添加成功")
	utils.RespSuccess(ctx, "添加用户成功", nil)
}

// 删除用户
func (u *user) Delete(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if !ok {
		utils.RespFalured(ctx, "获取接口参数失败", nil)
	}
	if !dto.DelUserFromDb(name) {
		utils.RespFalured(ctx, "删除用户失败", nil)
		return
	}
	utils.Logger.Info("用户: ", name, "删除成功")
	utils.RespSuccess(ctx, "删除用户成功", nil)
}

// 更新用户
func (u *user) Update(ctx *gin.Context) {
	if err := ctx.BindJSON(&u.MUUpt); err != nil {
		utils.RespFalured(ctx, "更新用户api接口数据失败", nil)
		return
	}
	if !dto.UptUserToDb(u.MUUpt.Name, u.MUUpt.User) {
		utils.RespFalured(ctx, "更新用户失败", nil)
		return
	}
	utils.Logger.Info("用户: ", u.MUUpt.Name, "更新成功")
	utils.RespSuccess(ctx, "更新用户成功", nil)
}

// 用户列表
func (u *user) Select(ctx *gin.Context) {
	user, err := dto.SelectUserFromDb()
	if err != nil {
		utils.RespFalured(ctx, "获取用户列表失败", nil)
	}
	utils.RespSuccess(ctx, "查询用户成功", user)
}

// 检查用户
func (u *user) Check(ctx *gin.Context) {
	name, _ := ctx.Get("name")
	user, err := dto.GetUserFromDb(name.(string))
	if err != nil {
		utils.RespFalured(ctx, "获取单个用户信息失败", nil)
	}
	utils.RespSuccess(ctx, "校验数据成功", &user)
}
