package auth

import "github.com/gin-gonic/gin"

func RouteAuth(r *gin.Engine) *gin.Engine {
	v1 := r.Group("v1")
	{
		// 系统
		v1.POST("/login", handleLogin)
		//	v1.GET("/userinfo", utils.MWtoken(), utils.MWUserCheck(), handleUserinfo)
		//	v1.GET("/menu", utils.MWtoken(), utils.MWUserCheck(), handleMenu)
	}

	return r
}
