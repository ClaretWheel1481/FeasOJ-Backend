package gincontext

import (
	"net/http"
	"src/internal/utils/sql"

	"github.com/gin-gonic/gin"
)

// 用户获取竞赛列表
func GetCompetitionList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": 200, "contests": sql.SelectCompetitionInfo()})
}
