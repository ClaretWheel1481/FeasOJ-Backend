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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 注册
func Register(c *gin.Context) {
	var req global.RegisterRequest
	c.ShouldBind(&req)
	// 判断用户或邮箱是否存在
	if sql.IsUserExist(req.Username, req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user or email already exists"})
		return
	}
	vcodeStatus := utils.CompareVerifyCode(req.Vcode, req.Email)
	if !vcodeStatus {
		c.JSON(http.StatusBadRequest, gin.H{"message": "captcha verification failed"})
		return
	}
	regstatus := sql.Register(req.Username, utils.EncryptPassword(req.Password), req.Email, uuid.New().String(), 0)
	if regstatus {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
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
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	} else {
		// 用户是否被封禁
		if sql.SelectUserInfo(Username).IsBan {
			c.JSON(http.StatusForbidden, gin.H{"message": "user is banned"})
			return
		}
		// 校验密码是否正确
		if utils.VerifyPassword(Password, userPassword) {
			// 生成Token并返回至前端
			token, err := utils.GenerateToken(Username)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "token generation failed"})
			}
			c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "password error"})
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
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
	} else {
		if sql.SelectUserByEmail(emails).Username != "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "user already exists"})
			return
		}
	}
	if utils.SendVerifycode(config.InitEmailConfig(), emails, utils.GenerateVerifycode()) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed, please try again"})
	}
}

// 验证用户信息
func VerifyUserInfo(c *gin.Context) {
	var Username string
	user := c.GetHeader("username")
	unescapeUsername, err := url.QueryUnescape(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't get username"})
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
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Info": userInfo})
	}
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 获取用户名
	username := c.Param("username")
	// 查询对应的用户信息
	userInfo := sql.SelectUserInfo(username)
	if userInfo.Username == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Info": userInfo})
	}
}

// 更新密码
func UpdatePassword(c *gin.Context) {
	var req global.UpdatePasswordRequest
	c.ShouldBind(&req)
	vcodeStatus := utils.CompareVerifyCode(req.Vcode, req.Email)
	if !vcodeStatus {
		c.JSON(http.StatusBadRequest, gin.H{"message": "captcha verification failed"})
		return
	}
	newPassword := utils.EncryptPassword(req.NewPassword)
	if sql.UpdatePassword(req.Email, newPassword) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
	}
}

// 更新个人简介
func UpdateSynopsis(c *gin.Context) {
	synopsis := c.PostForm("synopsis")
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	// 更新简介
	if sql.UpdateSynopsis(username, synopsis) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed"})
	}
}

// 上传头像
func UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't get file"})
		return
	}
	encodedUsername := c.GetHeader("username")
	username, _ := url.QueryUnescape(encodedUsername)
	// 获取用户信息
	userInfo := sql.SelectUserInfo(username)
	if userInfo.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't get user info"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	// 压缩图像
	if err := utils.CompressImage(originalFilePath, compressedFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	// 上传压缩后的头像路径至数据库
	if !sql.UpdateAvatar(username, newFilename) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
