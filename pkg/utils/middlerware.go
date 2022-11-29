package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// jwt
func JwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "token不能为空"})
			ctx.Abort()
			return
		}
		if claims, err := ParseToken(token); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "token解析失败"})
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
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "您的ip禁止访问该应用"})
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

// API接口白名单
func ApiWhite() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rPath := ctx.Request.RequestURI
		token := ctx.Request.Header.Get("token")
		if !isInSilence(rPath, viper.GetStringSlice("app_white_api")) {
			if token == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "token不能为空"})
				ctx.Abort()
				return
			} else {
				if claims, err := ParseToken(token); err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "token解析失败"})
					ctx.Abort()
					return
				} else {
					ctx.Set("name", claims.Username)
					ctx.Next()
				}
			}
		} else {
			ctx.Next()
		}
	}
}
