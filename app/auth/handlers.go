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

func handleUserInfoGet(c *gin.Context) {
	username, _ := c.Get("username")
	config, err := utils.ReadConfig()
	if err != nil {
		utils.Logger.Error("读取配置文件失败", "err", err)
		utils.RespondWithError(c, 500)
		return
	}

	user := config.GetUserWithUsername(username.(string))

	// user 必然存在 由中间件检查
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "success", "data": gin.H{
		"username": username,
		"nickname": user.Nickname,
		"role":     user.Role,
		"avatar":   user.Avatar,
	}})
}

func handleMenuGet(c *gin.Context) {
	type menuItem struct {
		ID        int        `json:"id"`
		Type      string     `json:"type"`
		Section   string     `json:"section"`
		Title     string     `json:"title"`
		Component string     `json:"component"`
		Icon      string     `json:"icon"`
		To        string     `json:"to"`
		Links     []menuItem `json:"links"`
	}

	var menu []menuItem

	home := menuItem{
		ID:        1,
		Type:      "link",
		Section:   "",
		Title:     "首页",
		Component: "home/index",
		Icon:      "ri-table-alt-line",
		To:        "/home",
		Links:     nil,
	}

	settings := menuItem{
		ID:        2,
		Type:      "group",
		Section:   "",
		Title:     "设置",
		Component: "",
		Icon:      "ri-home-smile-line",
		To:        "",
		Links: []menuItem{
			{
				ID:        20001,
				Type:      "link",
				Section:   "",
				Title:     "Player",
				Component: "settings/player",
				Icon:      "ri-home-smile-line",
				To:        "/settings/player",
				Links:     nil,
			},
			{
				ID:        20002,
				Type:      "link",
				Section:   "",
				Title:     "Room",
				Component: "settings/room",
				Icon:      "ri-home-smile-line",
				To:        "/settings/room",
				Links:     nil,
			},
		},
	}

	menu = append(menu, home, settings)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    menu,
	})
}

func handleRegisterPost(c *gin.Context) {
	if utils.Registered {
		utils.RespondWithError(c, 425)
		return
	}

	var user utils.User
	if err := c.ShouldBindJSON(&user); err != nil {
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

	user.Role = "admin"
	user.Disabled = false
	config.Users = append(config.Users, user)
	config.Registered = true
	utils.Registered = true

	err = utils.WriteConfig(config)
	if err != nil {
		utils.Logger.Error("写入配置文件失败", "err", err)
		utils.RespondWithError(c, 500)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data":    nil,
	})
}

func handleRegisterGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    utils.Registered,
	})
}
