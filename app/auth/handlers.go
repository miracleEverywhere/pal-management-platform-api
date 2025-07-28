package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pal-management-platform-api/utils"
)

func handleLogin(c *gin.Context) {
	type LoginForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginForm LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		// 如果绑定失败，返回 400 错误
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("读取配置文件失败", "err", err)
		utils.RespondWithError(c, 500)
		return
	}
	// 校验用户名和密码
	for _, user := range config.Users {
		if loginForm.Username == user.Username {
			if user.Disabled {
				utils.RespondWithError(c, 423)
				return
			}
			if loginForm.Password == user.Password {
				jwtSecret := []byte(config.JwtSecret)
				token, _ := utils.GenerateJWT(user, jwtSecret, 12)
				c.JSON(http.StatusOK, gin.H{"code": 200, "message": "登录成功", "data": gin.H{"token": token}})
				return
			} else {
				utils.RespondWithError(c, 422)
				return
			}
		}
	}

	utils.RespondWithError(c, 421)
}
