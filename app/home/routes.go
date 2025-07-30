package home

import (
	"github.com/gin-gonic/gin"
	"pal-management-platform-api/utils"
)

func RouteHome(r *gin.Engine) *gin.Engine {
	v1 := r.Group("v1")
	v1.Use(utils.MWtoken())
	v1.Use(utils.MWUserCheck())
	{
		home := v1.Group("home")
		{
			home.GET("/sys_info", handleSystemInfoGet)
		}
	}

	return r
}
