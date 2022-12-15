package route

import (
	"GThree/pkg/route/api"
	"GThree/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetRouter() *gin.Engine {
	route := gin.New()
	// app运行模式
	if viper.GetBool("app_debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// 添加全局中间件
	route.Use(utils.IpWhite())
	// 用户接口
	apiUser := api.Newuser()
	user := route.Group("/api")
	{
		user.POST("/sign", apiUser.Sign)
		user.Use(utils.JwtAuth())
		user.POST("/user", apiUser.Add)
		user.DELETE("/user/:name", apiUser.Delete)
		user.PUT("/user", apiUser.Update)
		user.GET("/users", apiUser.Select)
		user.GET("/user", apiUser.Check)
	}
	// 区服接口
	apiZone := api.NewZone()
	zone := route.Group("/api")
	{
		zone.Use(utils.JwtAuth())
		zone.GET("/zone/:name", apiZone.Result)
		zone.POST("/zone", apiZone.Manage)
	}

	return route
}
