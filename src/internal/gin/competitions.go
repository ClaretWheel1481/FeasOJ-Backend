package gincontext

import (
	"net/http"
	"net/url"
	"src/internal/utils/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 用户获取竞赛列表
func GetCompetitionList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contests": sql.SelectCompetitionInfo()})
}

// 用户加入竞赛
func JoinCompetition(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	uid := sql.SelectUserInfo(username).Uid
	if sql.AddUserCompetition(uid, competitionIdInt) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
}

// 查询用户是否在竞赛中
func IsInCompetition(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	uid := sql.SelectUserInfo(username).Uid
	if sql.SelectUserCompetition(uid, competitionIdInt) {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "isIn": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success", "isIn": false})
}

// 查询指定竞赛中的所有参与用户
func GetCompetitionUsers(c *gin.Context) {
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	c.JSON(http.StatusOK, gin.H{"users": sql.SelectUsersCompetition(competitionIdInt)})
}

// 用户退出竞赛
func QuitCompetition(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	uid := sql.SelectUserInfo(username).Uid
	if sql.DeleteUserCompetition(uid, competitionIdInt) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
}
