package gincontext

import (
	"net/http"
	"net/url"
	"src/internal/global"
	"src/internal/utils/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 管理员获取所有题目
func GetAllProblemsAdmin(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	problems := sql.SelectAllProblemsAdmin()
	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

// 管理员获取指定题目所有信息
func GetProblemAllInfo(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	problemInfo := sql.SelectProblemTestCases(c.Param("Pid"))
	c.JSON(http.StatusOK, gin.H{"problemInfo": problemInfo})
}

// 更新题目信息
func UpdateProblemInfo(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	var req global.AdminProblemInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}

	// 更新题目信息
	if err := sql.UpdateProblem(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 删除题目及其输入输出样例
func DeleteProblem(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	pid := c.Param("Pid")
	pidInt, _ := strconv.Atoi(pid)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if !sql.DeleteProblemAllInfo(pidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 管理员获取所有用户信息
func GetAllUsersInfo(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)

	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	usersInfo := sql.SelectAllUsersInfo()
	c.JSON(http.StatusOK, gin.H{"usersInfo": usersInfo})
}

// 晋升用户
func PromoteUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if !sql.PromoteToAdmin(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 降级用户
func DemoteUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if !sql.DemoteToUser(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 封禁用户
func BanUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if !sql.BanUser(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 解封用户
func UnbanUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if !sql.UnbanUser(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 管理员获取竞赛列表
func GetCompetitionListAdmin(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contests": sql.SelectCompetitionInfoAdmin()})
}

// 管理员获取指定竞赛ID信息
func GetCompetitionInfoAdmin(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	cid := c.Param("cid")
	cidInt, _ := strconv.Atoi(cid)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"contest": sql.SelectCompetitionInfoAdminByCid(cidInt)})
}

// 删除指定ID竞赛
func DeleteCompetition(c *gin.Context) {
	cid := c.Param("cid")
	cidInt, _ := strconv.Atoi(cid)
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}
	if !sql.DeleteCompetition(cidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// 更新/添加竞赛信息
func UpdateCompetitionInfo(c *gin.Context) {
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	var req global.AdminCompetitionInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if sql.SelectUserInfo(username).Role != 1 {
		c.JSON(http.StatusForbidden, gin.H{"message": "permission denied"})
		return
	}

	// 更新题目信息
	if err := sql.UpdateCompetition(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
