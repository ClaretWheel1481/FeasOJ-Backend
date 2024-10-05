package middleware

import (
	"net/http"
	"net/url"
	"src/utils"
	"src/utils/sql"

	"github.com/gin-gonic/gin"
)

func HeaderVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var User string
		encodedUsername := c.GetHeader("username")
		username, err := url.QueryUnescape(encodedUsername)
		token := c.GetHeader("Authorization")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user not found"})
			c.Abort()
			return
		}
		if utils.IsEmail(username) {
			User = sql.SelectUserByEmail(username).Username
		} else {
			User = username
		}
		if !utils.VerifyToken(User, token) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "can't not verify token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
