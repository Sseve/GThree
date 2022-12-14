package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ret(ctx *gin.Context, status, code int, msg string, data interface{}) {
	ctx.JSON(status, gin.H{"code": code, "message": msg, "data": data})
}

func RespSuccess(ctx *gin.Context, msg string, data interface{}) {
	ret(ctx, http.StatusOK, 20000, msg, data)
}

func RespFalured(ctx *gin.Context, msg string, data interface{}) {
	ret(ctx, http.StatusBadRequest, 40000, msg, data)
}
