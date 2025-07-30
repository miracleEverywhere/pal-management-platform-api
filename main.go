package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pal-management-platform-api/app/auth"
	"pal-management-platform-api/app/home"
	"pal-management-platform-api/utils"
	"runtime"
)

func main() {
	if utils.VersionShow {
		fmt.Println(utils.VERSION + "\n" + runtime.Version())
		return
	}

	gin.DefaultWriter = utils.InitAccessLogger()
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.Default()

	r = auth.RouteAuth(r)
	r = home.RouteHome(r)

	err := r.Run(fmt.Sprintf(":%d", utils.BindPort))
	if err != nil {
		utils.Logger.Error("启动服务器失败", "err", err)
		panic(err)
	}
}

func init() {
	// 绑定flag
	utils.BindFlags()
	// 数据库检查
	utils.CheckConfig()
	// 设置全局变量
	utils.SetGlobalVariables()
	// 检查目录
	utils.CheckDirs()
}
