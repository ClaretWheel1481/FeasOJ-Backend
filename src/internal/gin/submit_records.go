package gincontext

import (
	"net/http"
	"src/internal/utils/sql"

	"github.com/gin-gonic/gin"
)

// 获取所有提交记录
func GetAllSubmitRecords(c *gin.Context) {
	submitrecords := sql.SelectAllSubmitRecords()
	c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
}

// 获取指定用户提交记录
func GetSubmitRecordsByUsername(c *gin.Context) {
	username := c.Param("username")
	uid := sql.SelectUserInfo(username).Uid
	submitrecords := sql.SelectSubmitRecordsByUid(uid)
	c.JSON(http.StatusOK, gin.H{"submitrecords": submitrecords})
}
