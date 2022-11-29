package api

import (
	"GThree/pkg/grpc/gtmaster"
	"GThree/pkg/models"
	"GThree/pkg/utils"
	"fmt"
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
		utils.Falured(ctx, "获取区服接口数据失败", nil)
		return
	}
	// 数据入库

	// 远程调用
	var wg sync.WaitGroup
	for _, zone := range z.MZOpt.Zone {
		wg.Add(1)
		go gtmaster.CallServant(zone)
	}
	data := make([]gtmaster.ZoneResponse, 0, len(z.MZOpt.Zone))
	go func() {
		for {
			zoneResp := <-gtmaster.ZoneResult
			data = append(data, zoneResp)
			wg.Done()
		}
	}()
	wg.Wait()
	fmt.Println(data)
	// 成功返回
	utils.Success(ctx, "区服操作成功", data)
}
