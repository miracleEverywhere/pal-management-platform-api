package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ValidateJWT(tokenString string, jwtSecret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		Logger.Warn("JWT验证失败")
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func MWtoken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")
		config, err := ReadConfig()
		if err != nil {
			Logger.Error("配置文件打开失败", "err", err)
			RespondWithError(c, 500)
			c.Abort()
			return
		}
		tokenSecret := config.JwtSecret
		claims, err := ValidateJWT(token, []byte(tokenSecret))
		if err != nil {
			RespondWithError(c, 420)
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Set("nickname", claims.Nickname)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// MWUserCheck 用户状态检查
func MWUserCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exist := c.Get("username")
		if exist {
			usernameStr, ok := username.(string)
			if ok {
				config, err := ReadConfig()
				if err != nil {
					Logger.Error("读取配置文件失败", "err", err)
					RespondWithError(c, 500)
					c.Abort()
					return
				}

				user, _ := config.GetUserWithUsername(usernameStr)
				if len(user.Username) != 0 {
					if !user.Disabled {
						c.Next()
						return
					} else {
						RespondWithError(c, 423)
						c.Abort()
						return
					}
				} else {
					RespondWithError(c, 421)
					c.Abort()
					return
				}
			}

		}

		RespondWithError(c, 500)
		c.Abort()
		return
	}
}
