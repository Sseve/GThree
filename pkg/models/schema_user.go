package models

// 用户登录接口数据模型
type UserSign struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 添加用户接口数据模型
type UserAdd struct {
	Name    string   `json:"name" binding:"required"`
	PassOne string   `json:"pass_one" binding:"required"`
	PassTow string   `json:"pass_tow" binding:"required"`
	Desc    string   `json:"desc" binding:"required"`
	Roles   []string `json:"roles"`
	Avatar  string   `json:"avatar"`
}
