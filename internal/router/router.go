package router

import (
	gincontext "src/internal/gin"
	"src/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadRouter(r *gin.Engine) *gin.RouterGroup {
	r.Use(middlewares.Logger())
	// 设置路由组
	router1 := r.Group("/api/v1")
	{
		// 注册
		router1.POST("/register", gincontext.Register)

		// 登录
		router1.GET("/login", gincontext.Login)

		// 获取验证码
		router1.GET("/captcha", gincontext.GetCaptcha)

		// 获取用户信息
		router1.GET("/users/:username", gincontext.GetUserInfo)

		// 获取指定用户的提交记录
		router1.GET("/users/:username/submitrecords", gincontext.GetSubmitRecordsByUsername)

		// 密码修改
		router1.PUT("/users/password", gincontext.UpdatePassword)

		// 通知
		router1.GET("/notification/:uid", gincontext.SSEHandler)
	}

	authGroup := router1.Group("")
	authGroup.Use(middlewares.HeaderVerify())
	{
		// 验证用户信息
		authGroup.GET("/verify", gincontext.VerifyUserInfo)

		// 获取指定帖子的评论
		authGroup.GET("/discussions/comments/:did", gincontext.GetComment)

		// 获取竞赛列表
		authGroup.GET("/competitions", gincontext.GetCompetitionsList)

		// 获取指定竞赛ID信息
		authGroup.GET("/competitions/info/:cid", gincontext.GetCompetitionInfoByID)

		// 获取竞赛参与的用户列表
		authGroup.GET("/competitions/info/:cid/users", gincontext.GetCompetitionUsers)

		// 获取指定竞赛的所有题目
		authGroup.GET("/competitions/info/:cid/problems", gincontext.GetProblemsByCompetitionID)

		// 获取用户是否在竞赛中
		authGroup.GET("/competitions/:cid/in", gincontext.IsInCompetition)

		// 获取所有题目
		authGroup.GET("/problems", gincontext.GetAllProblems)

		// 获取所有帖子
		authGroup.GET("/discussions", gincontext.GetAllDiscussions)

		// 获取排行榜
		authGroup.GET("/ranking", gincontext.GetRanking)

		// 获取指定题目ID的所有信息
		authGroup.GET("/problems/:id", gincontext.GetProblemInfo)

		// 获取总提交记录
		authGroup.GET("/submitrecords", gincontext.GetAllSubmitRecords)

		// 获取指定帖子
		authGroup.GET("/discussions/:did", gincontext.GetDiscussionByDid)

		// 上传代码
		authGroup.POST("/problems/:pid/code", gincontext.UploadCode)

		// 创建讨论
		authGroup.POST("/discussions", gincontext.CreateDiscussion)

		// 添加评论
		authGroup.POST("/discussions/comments/:did", gincontext.AddComment)

		// 加入有密码的竞赛
		authGroup.POST("/competitions/join/pwd/:cid", gincontext.JoinCompetitionWithPassword)

		// 加入竞赛
		authGroup.POST("/competitions/join/:cid", gincontext.JoinCompetition)

		// 退出竞赛
		authGroup.POST("/competitions/quit/:cid", gincontext.QuitCompetition)

		// 用户上传头像
		authGroup.PUT("/users/avatar", gincontext.UploadAvatar)

		// 简介更新
		authGroup.PUT("/users/synopsis", gincontext.UpdateSynopsis)

		// 删除讨论
		authGroup.DELETE("/discussions/:did", gincontext.DeleteDiscussion)

		// 删除评论
		authGroup.DELETE("/discussions/comments/:cid", gincontext.DelComment)

	}

	adminGroup := authGroup.Group("/admin")
	// 管理员权限检查
	adminGroup.Use(middlewares.PermissionChecker())
	{
		// 管理员晋升用户
		adminGroup.PUT("/users/promote", gincontext.PromoteUser)

		// 管理员降级用户
		adminGroup.PUT("/users/demote", gincontext.DemoteUser)

		// 管理员封禁用户
		adminGroup.PUT("/users/ban", gincontext.BanUser)

		// 管理员解封用户
		adminGroup.PUT("/users/unban", gincontext.UnbanUser)

		// 管理员获取竞赛列表
		adminGroup.GET("/competitions", gincontext.GetCompetitionListAdmin)

		// 管理员获取所有题目
		adminGroup.GET("/problems", gincontext.GetAllProblemsAdmin)

		// 管理员获取指定竞赛ID信息
		adminGroup.GET("/competitions/:cid", gincontext.GetCompetitionInfoAdmin)

		// 管理员获取指定题目的所有信息
		adminGroup.GET("/problems/:pid", gincontext.GetProblemAllInfo)

		// 管理员获取所有用户信息
		adminGroup.GET("/users", gincontext.GetAllUsersInfo)

		// 管理员新增/更新题目信息
		adminGroup.POST("/problems", gincontext.UpdateProblemInfo)

		// 管理员新增/更新竞赛信息
		adminGroup.POST("/competitions", gincontext.UpdateCompetitionInfo)

		// 管理员删除题目
		adminGroup.DELETE("/problems/:pid", gincontext.DeleteProblem)

		// 管理员删除竞赛
		adminGroup.DELETE("/competitions/:cid", gincontext.DeleteCompetition)

		// 管理员启用竞赛计分
		adminGroup.GET("/competitions/:cid/score", gincontext.CalculateScore)

		// 管理员查看竞赛得分情况
		adminGroup.GET("/competitions/:cid/scoreboard", gincontext.GetScoreBoard)
	}
	return router1
}
