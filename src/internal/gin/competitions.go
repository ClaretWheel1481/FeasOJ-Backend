package gincontext

import (
	"net/http"
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
	uid := c.Param("uid")
	uidInt, _ := strconv.Atoi(uid)
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	if sql.AddUserCompetition(uidInt, competitionIdInt) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
}

// 查询用户是否在竞赛中
func IsInCompetition(c *gin.Context) {
	uid := c.Param("uid")
	uidInt, _ := strconv.Atoi(uid)
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	if sql.SelectUserCompetition(uidInt, competitionIdInt) {
		c.JSON(http.StatusOK, gin.H{"message": "Success", "isIn": true})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success", "isIn": false})
}

// 用户退出竞赛
func QuitCompetition(c *gin.Context) {
	uid := c.Param("uid")
	uidInt, _ := strconv.Atoi(uid)
	competitionId := c.Param("cid")
	competitionIdInt, _ := strconv.Atoi(competitionId)
	if sql.DeleteUserCompetition(uidInt, competitionIdInt) == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Fail"})
}
