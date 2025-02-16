package gincontext

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"src/config"
	"src/internal/global"
	"src/internal/utils"
	"src/internal/utils/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 注册
func Register(c *gin.Context) {
	var req global.RegisterRequest
	c.ShouldBind(&req)
	// 判断用户或邮箱是否存在
	if sql.IsUserExist(req.Username, req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "userAlreadyinUse")})
		return
	}
	if utils.ContainsProfanity(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
		return
	}
	vcodeStatus := utils.CompareVerifyCode(req.Vcode, req.Email)
	if !vcodeStatus {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "captchaError")})
		return
	}
	regstatus := sql.Register(req.Username, utils.EncryptPassword(req.Password), req.Email, uuid.New().String(), 0)
	if regstatus {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
}

// 登录
func Login(c *gin.Context) {
	var Username string
	user := c.Query("username")
	Password := c.Query("password")
	// 判断用户输入是否为邮箱
	if utils.IsEmail(user) {
		Username = sql.SelectUserByEmail(user).Username
	} else {
		Username = user
	}
	userPassword := utils.SelectUser(Username).Password
	if userPassword == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	} else {
		// 用户是否被封禁
		if sql.SelectUserInfo(Username).IsBan {
			c.JSON(http.StatusForbidden, gin.H{"message": GetMessage(c, "userIsBanned")})
			return
		}
		// 校验密码是否正确
		if utils.VerifyPassword(Password, userPassword) {
			// 生成Token并返回至前端
			token, err := utils.GenerateToken(Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			}
			c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "loginSuccess"), "token": token})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "wrongPassword")})
		}
	}
}

// 获取验证码
func GetCaptcha(c *gin.Context) {
	// 获取邮箱地址
	emails := c.Query("email")
	isCreate := c.GetHeader("iscreate")
	if isCreate == "false" {
		if sql.SelectUserByEmail(emails).Username == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}
	} else {
		if sql.SelectUserByEmail(emails).Username != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "userAlreadyinUse")})
			return
		}
	}
	if utils.SendVerifycode(config.InitEmailConfig(), emails, utils.GenerateVerifycode()) {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	}
}

// 验证用户信息
func VerifyUserInfo(c *gin.Context) {
	var Username string
	user := c.GetHeader("Username")
	unescapeUsername, err := url.QueryUnescape(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	if utils.IsEmail(unescapeUsername) {
		Username = sql.SelectUserByEmail(unescapeUsername).Username
	} else {
		Username = unescapeUsername
	}
	// 查询对应的用户信息
	userInfo := sql.SelectUserInfo(Username)
	if userInfo.Username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	} else {
		c.JSON(http.StatusOK, gin.H{"info": userInfo})
	}
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 获取用户名
	username := c.Param("username")
	// 查询对应的用户信息
	userInfo := sql.SelectUserInfo(username)
	if userInfo.Username == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
	} else {
		c.JSON(http.StatusOK, gin.H{"info": userInfo})
	}
}

// 更新密码
func UpdatePassword(c *gin.Context) {
	var req global.UpdatePasswordRequest
	c.ShouldBind(&req)
	vcodeStatus := utils.CompareVerifyCode(req.Vcode, req.Email)
	if !vcodeStatus {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "captchaError")})
		return
	}
	newPassword := utils.EncryptPassword(req.NewPassword)
	if sql.UpdatePassword(req.Email, newPassword) {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
}

// 更新个人简介
func UpdateSynopsis(c *gin.Context) {
	synopsis := c.PostForm("synopsis")
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	// 更新简介
	if utils.ContainsProfanity(synopsis) {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "profanity")})
		return
	}
	if sql.UpdateSynopsis(username, synopsis) {
		c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "failed")})
	}
}

// 上传头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}
	encodedUsername := c.GetHeader("Username")
	username, _ := url.QueryUnescape(encodedUsername)
	// 获取用户信息
	userInfo := sql.SelectUserInfo(username)
	if userInfo.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": GetMessage(c, "invalidrequest")})
		return
	}
	newFilename := fmt.Sprintf("%d%s", userInfo.Uid, path.Ext(file.Filename))
	originalFilePath := filepath.Join(global.AvatarsDir, newFilename)
	compressedFilePath := filepath.Join(global.AvatarsDir, fmt.Sprintf("%d%s", userInfo.Uid, path.Ext(file.Filename)))
	if _, err := os.Stat(compressedFilePath); err == nil {
		os.Remove(compressedFilePath)
	}
	// 保存原始文件
	if err := c.SaveUploadedFile(file, originalFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	// 压缩图像
	if err := utils.CompressImage(originalFilePath, compressedFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	// 上传压缩后的头像路径至数据库
	if !sql.UpdateAvatar(username, newFilename) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": GetMessage(c, "success")})
}

// 获取排行榜
func GetRanking(c *gin.Context) {
	cacheKey := "ranking"
	// 从缓存中获取数据
	var ranking []global.UserInfoRequest
	err := utils.GetCache(cacheKey, &ranking)
	if err != nil || len(ranking) == 0 {
		ranking = sql.SelectRank100Users()
		// 缓存5分钟
		err := utils.SetCache(cacheKey, ranking, 5*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": GetMessage(c, "internalServerError")})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"ranking": ranking})
}
