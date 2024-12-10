package router

import (
	gincontext "src/internal/gin"
	"src/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func LoadRouter(r *gin.Engine) *gin.RouterGroup {
	// 设置路由组
	router1 := r.Group("/api/v1")
	router1.Use(middlewares.Logger())
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

		// 获取竞赛参与的用户列表
		router1.GET("/competitions/:cid/users", gincontext.GetCompetitionUsers)

		router1.Use(middlewares.HeaderVerify())
		{
			// 验证用户信息
			router1.GET("/verify", gincontext.VerifyUserInfo)

			// 获取指定帖子的评论
			router1.GET("/discussions/comments/:did", gincontext.GetComment)

			// 获取竞赛列表
			router1.GET("/competitions", gincontext.GetCompetitionList)

			// 获取用户是否在竞赛中
			router1.GET("/competitions/:cid/in", gincontext.IsInCompetition)

			// 获取所有题目
			router1.GET("/problems", gincontext.GetAllProblems)

			// 获取所有讨论帖子
			router1.GET("/discussions", gincontext.GetAllDiscussions)

			// 根据题目ID获取题目信息
			router1.GET("/problems/:id", gincontext.GetProblemInfo)

			// 获取总提交记录
			router1.GET("/submitrecords", gincontext.GetAllSubmitRecords)

			// 获取指定帖子
			router1.GET("/discussions/:did", gincontext.GetDiscussionByDid)

			// 上传代码
			router1.POST("/problems/:pid/code", gincontext.UploadCode)

			// 创建讨论
			router1.POST("/discussions", gincontext.CreateDiscussion)

			// 添加评论
			router1.POST("/discussions/comments/:did", gincontext.AddComment)

			// 加入竞赛
			router1.POST("/competitions/join/:cid", gincontext.JoinCompetition)

			// 退出竞赛
			router1.POST("/competitions/quit/:cid", gincontext.QuitCompetition)

			// 用户上传头像
			router1.PUT("/users/avatar", gincontext.UploadAvatar)

			// 简介更新
			router1.PUT("/users/synopsis", gincontext.UpdateSynopsis)

			// 删除讨论
			router1.DELETE("/discussions/:did", gincontext.DeleteDiscussion)

			// 删除评论
			router1.DELETE("/discussions/comments/:cid", gincontext.DelComment)

			// 管理员权限检查
			router1.Use(middlewares.PermissionChecker())
			{
				// 管理员晋升用户
				router1.PUT("/admin/users/promote", gincontext.PromoteUser)

				// 管理员降级用户
				router1.PUT("/admin/users/demote", gincontext.DemoteUser)

				// 管理员封禁用户
				router1.PUT("/admin/users/ban", gincontext.BanUser)

				// 管理员解封用户
				router1.PUT("/admin/users/unban", gincontext.UnbanUser)

				// 管理员获取竞赛列表
				router1.GET("/admin/competitions", gincontext.GetCompetitionListAdmin)

				// 管理员获取所有题目
				router1.GET("/admin/problems", gincontext.GetAllProblemsAdmin)

				// 管理员获取指定竞赛ID信息
				router1.GET("/admin/competitions/:cid", gincontext.GetCompetitionInfoAdmin)

				// 管理员获取指定题目的所有信息
				router1.GET("/admin/problems/:pid", gincontext.GetProblemAllInfo)

				// 管理员获取所有用户信息
				router1.GET("/admin/users", gincontext.GetAllUsersInfo)

				// 管理员新增/更新题目信息
				router1.POST("/admin/problems", gincontext.UpdateProblemInfo)

				// 管理员新增/更新竞赛信息
				router1.POST("/admin/competitions", gincontext.UpdateCompetitionInfo)

				// 管理员删除题目
				router1.DELETE("/admin/problems/:pid", gincontext.DeleteProblem)

				// 管理员删除竞赛
				router1.DELETE("/admin/competitions/:cid", gincontext.DeleteCompetition)
			}
		}
	}
	return router1
}
