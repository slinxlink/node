package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 公开路由（不需要登录）
	public := r.Group("/api")
	{
		public.POST("/auth/login", Login)
	}

	// 需要登录的路由
	private := r.Group("/api")
	private.Use(AuthMiddleware())
	{
		// 流程
		private.POST("/restart", Restart)

		// 认证
		private.POST("/auth/logout", Logout)
		private.PUT("/auth", ChangeCredentials)

		// 配置
		private.GET("/config", GetConfig)
		private.PUT("/config", UpdateConfig)
		private.POST("/config/cert/cf", ApplyCertCF)
		private.POST("/config/cert/http", ApplyCertHTTP)

		// 入站
		private.GET("/inbound", GetInbounds)
		private.PUT("/inbound/save", SaveInbound)
		private.DELETE("/inbound/:id", DeleteInbound)
		private.POST("/inbound/quick", QuickInbound)
		private.PUT("/inbound/:id/toggle", ToggleInbound)

		// 用户
		private.GET("/user", GetUsers)
		private.PUT("/user/save", SaveUser)
		private.DELETE("/user/:id", DeleteUser)
		private.PUT("/user/:id/toggle", ToggleUser)

		// 对接
		private.GET("/board", GetBoards)
		private.PUT("/board/save", SaveBoard)
		private.DELETE("/board/:id", DeleteBoard)
		private.PUT("/board/:id/toggle", ToggleBoard)

		private.GET("/board/:id/user", GetBoardUser)

		// 生成
		private.GET("/generate/port", GeneratePort)
		private.GET("/generate/reality-target", GenerateRealityTarget)
		private.GET("/generate/reality-keypair", GenerateRealityKeyPair)
		private.GET("/generate/shortids", GenerateShortIDs)
		private.GET("/generate/token", GenerateToken)
		private.GET("/generate/uuid", GenerateUUID)
		private.GET("/generate/password", GeneratePassword)

		// 检测
		private.POST("/detect/ip/fetch", FetchIP)
		private.GET("/detect/ip", DetectIP)
		private.POST("/detect/unlock/fetch", FetchUnlock)
		private.GET("/detect/unlock", DetectUnlock)
		private.GET("/detect/route", DetectRoute)
		private.POST("/detect/route/fetch", FetchRoute)

		// 系统
		private.GET("/system/status", GetSystemStatus)
		private.GET("/stats", GetStats)
		private.GET("/system/log", GetSystemLog)

		// 核心
		private.GET("/core", GetCore)
		private.PUT("/core", UpdateCore)
		private.POST("/core/reset", ResetCore)
		private.GET("/core/status", GetCoreStatus)
		private.POST("/core/start", StartCore)
		private.POST("/core/stop", StopCore)
		private.POST("/core/restart", RestartCore)
		private.GET("/core/config", GetCoreConfig)
		private.GET("/core/process", GetCoreProcess)

		// 日志
		private.GET("/log/slinx", SlinxLog)
		private.GET("/log/core", CoreLog)

		// 推送
		private.GET("/task/:id", TaskLog)

		// 证书
		private.GET("/cert", GetCert)
		private.POST("/cert", SaveCert)
		private.DELETE("/cert/:id", DeleteCert)
		private.POST("/cert/:id/apply", ApplyCert)
		private.GET("/cert/:id/content", GetCertContent)

		private.GET("/acme", GetAcme)
		private.POST("/acme", SaveAcme)
		private.DELETE("/acme/:id", DeleteAcme)

		private.GET("/dns", GetDnsAccount)
		private.POST("/dns", SaveDnsAccount)
		private.DELETE("/dns/:id", DeleteDnsAccount)
	}

	// 订阅链接（公开）
	// r.GET("/sub/:token", GetSubscription)
}
