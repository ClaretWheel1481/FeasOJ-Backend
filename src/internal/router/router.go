package router

import (
	gincontext "src/internal/gin"
	"src/internal/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRouter(r *gin.Engine) *gin.RouterGroup {
	// 设置路由组
	router1 := r.Group("/api/v1")
	router1.Use(middleware.Logger())
	{
		// 注册API
		router1.POST("/register", gincontext.Register)

		// 登录API
		router1.GET("/login", gincontext.Login)

		// 获取验证码API
		router1.GET("/captcha", gincontext.GetCaptcha)

		// 获取用户信息API
		router1.GET("/users/:username", gincontext.GetUserInfo)

		// 获取指定用户的提交记录API
		router1.GET("/users/:username/submitrecords", gincontext.GetSubmitRecordsByUsername)

		// 密码修改API
		router1.PUT("/users/password", gincontext.UpdatePassword)

		router1.Use(middleware.HeaderVerify())
		{
			// 验证用户信息API
			router1.GET("/verify", gincontext.VerifyUserInfo)

			// 获取指定帖子的评论API
			router1.GET("/discussions/:Did/comments", gincontext.GetComment)

			// 获取竞赛列表
			router1.GET("/competitions", gincontext.GetCompetitionList)

			// 获取所有题目API
			router1.GET("/problems", gincontext.GetAllProblems)

			// 获取所有讨论帖子API
			router1.GET("/discussions", gincontext.GetAllDiscussions)

			// 根据题目ID获取题目信息
			router1.GET("/problems/:id", gincontext.GetProblemInfo)

			// 获取总提交记录API
			router1.GET("/submitrecords", gincontext.GetAllSubmitRecords)

			// 获取指定Did的帖子API
			router1.GET("/discussions/:Did", gincontext.GetDiscussionByDid)

			// 上传代码API
			router1.POST("/users/code", gincontext.UploadCode)

			// 创建讨论API
			router1.POST("/discussions", gincontext.CreateDiscussion)

			// 添加评论API
			router1.POST("/discussions/:Did/comments", gincontext.AddComment)

			// 用户上传头像API
			router1.PUT("/users/avatar", gincontext.UploadAvatar)

			// 简介更新API
			router1.PUT("/users/synopsis", gincontext.UpdateSynopsis)

			// 管理员晋升用户API
			router1.PUT("/admin/users/promote", gincontext.PromoteUser)

			// 管理员降级用户API
			router1.PUT("/admin/users/demote", gincontext.DemoteUser)

			// 管理员封禁用户API
			router1.PUT("/admin/users/ban", gincontext.BanUser)

			// 管理员解封用户API
			router1.PUT("/admin/users/unban", gincontext.UnbanUser)

			// 管理员获取竞赛列表
			router1.GET("/admin/competitions", gincontext.GetCompetitionListAdmin)

			// 管理员获取所有题目API
			router1.GET("/admin/problems", gincontext.GetAllProblemsAdmin)

			// 管理员获取指定竞赛ID信息
			router1.GET("/admin/competitions/:cid", gincontext.GetCompetitionInfoAdmin)

			// 管理员获取指定题目的所有信息API
			router1.GET("/admin/problems/:Pid", gincontext.GetProblemAllInfo)

			// 管理员获取所有用户信息API
			router1.GET("/admin/users", gincontext.GetAllUsersInfo)

			// 管理员新增/更新题目信息API
			router1.POST("/admin/problems", gincontext.UpdateProblemInfo)

			// 管理员新增/更新竞赛信息API
			router1.POST("/admin/competitions", gincontext.UpdateCompetitionInfo)

			// 管理员删除题目API
			router1.DELETE("/admin/problems/:Pid", gincontext.DeleteProblem)

			// 管理员删除竞赛API
			router1.DELETE("/admin/competitions/:cid", gincontext.DeleteCompetition)

			// 删除讨论API
			router1.DELETE("/discussions/:Did", gincontext.DeleteDiscussion)

			// 删除评论API
			router1.DELETE("/comments/:Cid", gincontext.DelComment)
		}
	}
	return router1
}
