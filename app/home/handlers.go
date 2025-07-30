package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pal-management-platform-api/utils"
)

func handleSystemInfoGet(c *gin.Context) {
	type Data struct {
		Cpu    float64 `json:"cpu"`
		Memory float64 `json:"memory"`
	}
	type Response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}

	var err error
	var SysInfoResponse Response
	SysInfoResponse.Code = 200
	SysInfoResponse.Message = "success"
	SysInfoResponse.Data.Cpu, err = utils.CpuUsage()
	if err != nil {
		utils.Logger.Error("获取Cpu使用率失败", "err", err)
	}
	SysInfoResponse.Data.Memory, err = utils.MemoryUsage()
	if err != nil {
		utils.Logger.Error("获取内存使用率失败", "err", err)
	}

	c.JSON(http.StatusOK, SysInfoResponse)
}
