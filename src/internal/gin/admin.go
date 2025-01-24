package gincontext

import (
	"gorm.io/gorm"
	"net/http"
	"src/internal/global"
	"src/internal/utils"
	"src/internal/utils/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 管理员获取所有题目
func GetAllProblemsAdmin(c *gin.Context) {
	problems := sql.SelectAllProblemsAdmin()
	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

// 管理员获取指定题目所有信息
func GetProblemAllInfo(c *gin.Context) {
	problemInfo := sql.SelectProblemTestCases(c.Param("pid"))
	c.JSON(http.StatusOK, gin.H{"problemInfo": problemInfo})
}

// 更新题目信息
func UpdateProblemInfo(c *gin.Context) {
	var req global.AdminProblemInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}

	// 更新题目信息
	if err := sql.UpdateProblem(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 删除题目及其输入输出样例
func DeleteProblem(c *gin.Context) {
	pid := c.Param("pid")
	pidInt, _ := strconv.Atoi(pid)
	if !sql.DeleteProblemAllInfo(pidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 管理员获取所有用户信息
func GetAllUsersInfo(c *gin.Context) {
	usersInfo := sql.SelectAllUsersInfo()
	c.JSON(http.StatusOK, gin.H{"usersInfo": usersInfo})
}

// 晋升用户
func PromoteUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if !sql.PromoteToAdmin(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 降级用户
func DemoteUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if sql.SelectAdminCount() <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}

	if !sql.DemoteToUser(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 封禁用户
func BanUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if !sql.BanUser(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 解封用户
func UnbanUser(c *gin.Context) {
	uid := c.Query("uid")
	uidInt, _ := strconv.Atoi(uid)

	if !sql.UnbanUser(uidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 管理员获取竞赛列表
func GetCompetitionListAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"contests": sql.SelectCompetitionInfoAdmin()})
}

// 管理员获取指定竞赛ID信息
func GetCompetitionInfoAdmin(c *gin.Context) {
	cid := c.Param("cid")
	cidInt, _ := strconv.Atoi(cid)
	c.JSON(http.StatusOK, gin.H{"contest": sql.SelectCompetitionInfoAdminByCid(cidInt)})
}

// 删除指定ID竞赛
func DeleteCompetition(c *gin.Context) {
	cid := c.Param("cid")
	cidInt, _ := strconv.Atoi(cid)

	if !sql.DeleteCompetition(cidInt) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 更新/添加竞赛信息
func UpdateCompetitionInfo(c *gin.Context) {
	var req global.AdminCompetitionInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}

	// 更新题目信息
	if err := sql.UpdateCompetition(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 计算分数
func CalculateScore(c *gin.Context) {
	competitionId, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid competition ID"})
		return
	}

	// 查询竞赛信息
	var competition global.Competition
	result := utils.ConnectSql().First(&competition, competitionId)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Competition not found"})
		return
	}

	// 查询竞赛参与用户
	users := sql.SelectUsersCompetition(competitionId)
	if len(users) == 0 {
		c.JSON(200, gin.H{"message": "No users in competition"})
		return
	}

	// 遍历所有参与竞赛的用户
	for _, user := range users {
		var submissions []global.SubmitRecord
		utils.ConnectSql().
			Where("uid = ? AND result = ? AND time BETWEEN ? AND ?",
				user.Uid,
				"Success",
				competition.Start_at,
				competition.End_at).
			Find(&submissions)

		// 计算分数
		score := 0
		for _, submission := range submissions {
			var difficulty string
			err := utils.ConnectSql().
				Table("problems").
				Select("difficulty").
				Where("contest_id = ? AND pid = ?", competitionId, submission.Pid).
				Row().
				Scan(&difficulty)
			if err != nil {
				continue
			}

			switch difficulty {
			case "简单":
				score += 1
			case "中等":
				score += 3
			case "困难":
				score += 5
			}
		}

		// 更新用户分数
		if score > 0 {
			utils.ConnectSql().Model(&global.User{}).Where("uid = ?", user.Uid).Update("score", gorm.Expr("score + ?", score))
		}

		// TODO: 额外建一个竞赛情况表，用于存储用户在某个竞赛中获取的分数，用于后台统计
	}
	c.JSON(200, gin.H{"message": "Scores calculated successfully"})
}
