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
}

func Newuser() *user {
	return new(user)
}

// 用户登录
func (u *user) Sign(ctx *gin.Context) {
	// 获取接口数据
	if err := ctx.BindJSON(&u.MUSign); err != nil {
		utils.Falured(ctx, "获取用户登录api接口数据失败", nil)
		return
	}
	// 校验数据
	if !dto.CheckUserFromDb(u.MUSign.Name, u.MUSign.Password) {
		utils.Falured(ctx, "用户名或密码错误", nil)
		return
	}
	// 创建token
	token, err := utils.CreateToken(u.MUSign.Name)
	if err != nil {
		utils.Falured(ctx, "创建token失败", nil)
		return
	}
	// 返回结果
	utils.Success(ctx, "登录成功", token)
}

// 添加用户
func (u *user) Add(ctx *gin.Context) {
	if err := ctx.BindJSON(&u.MUAdd); err != nil {
		utils.Falured(ctx, "获取添加用户api接口数据失败", nil)
		return
	}
	if !dto.AddUserToDb(u.MUAdd) {
		utils.Falured(ctx, "添加用户失败,或许用户已经存在", nil)
		return
	}
	utils.Success(ctx, "添加用户成功", nil)
}

// 删除用户
func (u *user) Delete(ctx *gin.Context) {
	name, ok := ctx.Params.Get("name")
	if !ok {
		utils.Falured(ctx, "获取接口参数失败", nil)
	}
	if !dto.DelUserFromDb(name) {
		utils.Falured(ctx, "删除用户失败", nil)
		return
	}
	utils.Success(ctx, "删除用户成功", nil)
}

// 更新用户
func (u *user) Update(ctx *gin.Context) {

}

// 查询用户
func (u *user) Select(ctx *gin.Context) {

}
