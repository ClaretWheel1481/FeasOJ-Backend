package gincontext

import (
	"fmt"
	"net/http"
	"net/url"
	"src/internal/utils"
	"src/internal/utils/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取所有讨论列表
func GetAllDiscussions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	itemsPerPage, _ := strconv.Atoi(c.DefaultQuery("itemsPerPage", "12"))

	discussions, total := sql.SelectDiscussList(page, itemsPerPage)
	c.JSON(http.StatusOK, gin.H{
		"discussions": discussions,
		"total":       total,
	})
}

// 获取指定id讨论信息
func GetDiscussionByDid(c *gin.Context) {
	did, _ := strconv.Atoi(c.Param("did"))
	discussion := sql.SelectDiscussionByDid(did)
	c.JSON(http.StatusOK, gin.H{"discussionInfo": discussion})
}

// 创建讨论
func CreateDiscussion(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	title := c.PostForm("title")
	content := c.PostForm("content")

	// 检测文本
	if utils.DetectText(title) || utils.DetectText(content) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
		return
	}

	// 获取用户ID
	userInfo := sql.SelectUserInfo(username)

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 设置频率限制键
	userRateLimitKey := fmt.Sprintf("discussionRateLimit:%d", userInfo.Uid)
	exists, _ := rdb.Exists(userRateLimitKey).Result()
	if exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 设置限流键
	rdb.Set(userRateLimitKey, 1, 15*time.Second)

	// 创建讨论
	uid := userInfo.Uid

	if !sql.AddDiscussion(title, content, uid) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 删除讨论
func DeleteDiscussion(c *gin.Context) {
	did, _ := strconv.Atoi(c.Param("did"))
	if sql.DelDiscussion(did) {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	}
}

// 获取指定讨论的评论
func GetComment(c *gin.Context) {
	did, _ := strconv.Atoi(c.Param("did"))
	comments := sql.SelectCommentsByDid(did)
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

// 删除指定Cid的评论
func DelComment(c *gin.Context) {
	cid, _ := strconv.Atoi(c.Param("cid"))
	if !sql.DeleteCommentByCid(cid) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 添加评论
func AddComment(c *gin.Context) {
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	content := c.PostForm("content")
	did, _ := strconv.Atoi(c.Param("did"))
	// 获取用户ID
	userInfo := sql.SelectUserInfo(username)

	rdb := utils.ConnectRedis()
	defer rdb.Close()

	// 设置频率限制键
	userRateLimitKey := fmt.Sprintf("commentRateLimit:%d", userInfo.Uid)
	exists, _ := rdb.Exists(userRateLimitKey).Result()
	if exists == 1 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": GetMessage(c, "rateLimit")})
		return
	}

	// 设置限流键
	rdb.Set(userRateLimitKey, 1, 10*time.Second)

	profanity := false
	if utils.DetectText(content) {
		profanity = true
	}
	if !sql.AddComment(content, did, userInfo.Uid, profanity) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}
