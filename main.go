package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pal-management-platform-api/store"
	"pal-management-platform-api/utils"
	"runtime"
)

func main() {
	if utils.VersionShow {
		fmt.Println(utils.VERSION + "\n" + runtime.Version())
		return
	}

	// 将 Gin 的默认输出重定向到文件
	gin.DefaultWriter = utils.InitAccessLogger()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})
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
	store.CheckConfig()
	// 设置全局变量
	utils.SetGlobalVariables()
	// 检查目录
	utils.CheckDirs()
}
