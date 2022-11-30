package api

import (
	"GThree/pkg/grpc/gtmaster"
	"GThree/pkg/models"
	"GThree/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type zone struct {
	MZOpt models.ZoneOpt
}

func NewZone() *zone {
	return new(zone)
}

// 区服管理
func (z *zone) Manage(ctx *gin.Context) {
	// 获取接口数据
	if err := ctx.BindJSON(&z.MZOpt); err != nil {
		utils.Falured(ctx, "获取区服接口数据失败", err)
		return
	}
	// 数据入库

	// 远程调用
	var wg sync.WaitGroup
	num := len(z.MZOpt.Zone)
	ZoneResult := make(chan gtmaster.ZoneResponse, num)
	wg.Add(num)
	for _, zone := range z.MZOpt.Zone {
		go gtmaster.CallServant(zone, ZoneResult)
	}
	data := make([]gtmaster.ZoneResponse, 0, num)
	go func() {
		for {
			data = append(data, <-ZoneResult)
			wg.Done()
		}
	}()
	wg.Wait()
	// 成功返回
	utils.Success(ctx, "区服操作成功", data)
}
