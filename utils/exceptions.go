package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var exceptions = map[int]string{
	404: "集群资源不存在",
	420: "Token认证失败",
	421: "用户不存在",
	422: "密码错误",
	423: "该用户已被禁用",
	424: "旧密码错误",
	425: "非法请求",
	429: "请求过于频繁，请稍后再试",
	500: "服务器内部错误",
	510: "获取主机信息失败",
	511: "执行命令失败",
}

func RespondWithError(c *gin.Context, code int) {
	message := exceptions[code]
	c.JSON(http.StatusOK, gin.H{"code": code, "message": message})
}
