package auth

import (
	"github.com/gin-gonic/gin"
	"pal-management-platform-api/utils"
)

func RouteAuth(r *gin.Engine) *gin.Engine {
	v1 := r.Group("v1")
	{
		// 系统
		v1.POST("/login", handleLogin)
		v1.GET("/userinfo", utils.MWtoken(), utils.MWUserCheck(), handleUserInfoGet)
		v1.GET("/menu", utils.MWtoken(), utils.MWUserCheck(), handleMenuGet)
		v1.POST("/register", handleRegisterPost)
		v1.GET("/register", handleRegisterGet)
	}

	return r
}
