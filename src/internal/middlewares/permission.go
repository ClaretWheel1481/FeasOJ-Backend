package middlewares

import (
	"net/http"
	"net/url"
	"src/internal/utils/sql"

	"github.com/gin-gonic/gin"
)

func PermissionChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		encodedUsername := c.GetHeader("Username")
		username, _ := url.QueryUnescape(encodedUsername)
		if sql.SelectUserInfo(username).Role != 1 {
			c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
