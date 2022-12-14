package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// jwt认证
func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token == "" {
			RespFalured(ctx, "token不能为空", nil)
			ctx.Abort()
			return
		}
		if claims, err := ParseToken(token); err != nil {
			RespFalured(ctx, "token解析失败", nil)
			ctx.Abort()
			return
		} else {
			ctx.Set("name", claims.Username)
			ctx.Next()
		}
	}
}

// ip白名单
func IpWhite() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rIp := ctx.RemoteIP()
		if !isInSilence(rIp, viper.GetStringSlice("app_white_ips")) {
			RespFalured(ctx, "您的IP禁止访问该应用", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func isInSilence(e string, es []string) bool {
	for _, v := range es {
		if v == e {
			return true
		}
	}
	return false
}
