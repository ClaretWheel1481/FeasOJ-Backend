package gincontext

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"src/internal/global"
	"src/internal/utils"
	"src/internal/utils/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取所有题目
func GetAllProblems(c *gin.Context) {
	// 实时性要求较高，不做数据缓存
	problems := sql.SelectAllProblems()
	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

// 获取题目信息
func GetProblemInfo(c *gin.Context) {
	// 生成缓存键
	cacheKey := "problemInfo_" + c.Param("id")
	var problemInfo global.ProblemInfoRequest

	// 从缓存中获取数据
	err := utils.GetCache(cacheKey, &problemInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	if problemInfo.Pid == 0 {
		// 缓存未命中
		problemInfo = sql.SelectProblemInfo(c.Param("id"))
		// 数据存入缓存，时间10分钟
		err = utils.SetCache(cacheKey, problemInfo, 10*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"problemInfo": problemInfo})
}

func UploadCode(c *gin.Context) {
	problem := c.Param("pid")
	pidInt, _ := strconv.Atoi(problem)
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	file, err := c.FormFile("code")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}

	// 获取用户ID
	userInfo := sql.SelectUserInfo(username)

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 提交频率限制
	userRateLimitKey := fmt.Sprintf("ratelimit:%d", userInfo.Uid)

	// 检查是否在限制时间内
	exists, _ := rdb.Exists(userRateLimitKey).Result()
	if exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 设置限流键，10秒后自动失效
	rdb.Set(userRateLimitKey, 1, 10*time.Second)

	// 将文件名改为用户ID_题目ID
	newFileName := fmt.Sprintf("%d_%s%s", userInfo.Uid, problem, path.Ext(file.Filename))
	filepath := filepath.Join(global.CodeDir, newFileName)

	// 保存文件到指定路径
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	var language string
	switch path.Ext(file.Filename) {
	case ".cpp":
		language = "C++"
	case ".java":
		language = "Java"
	case ".py":
		language = "Python"
	case ".go":
		language = "Go"
	default:
		language = "Unknown"
	}

	// 上传任务至Redis任务队列
	err = rdb.RPush("judgeTask", newFileName).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	sql.AddSubmitRecord(userInfo.Uid, pidInt, "Running...", language, username)
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}
