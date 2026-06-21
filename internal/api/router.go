package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 公开
	public := r.Group("/api")
	{
		public.POST("/auth/login", Login)
	}

	// 内部
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
		private.POST("/config/reset", ResetConfig)

		// 更新
		private.GET("/update/version", GetVersion)
		private.GET("/update/check", CheckUpdate)
		private.POST("/update", Update)

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

		// 订阅
		private.POST("/sub/uri", GetSubscriptionUri)
		private.POST("/sub/url", GetSubscriptionUrl)

		// 系统
		private.GET("/system/status", GetSystemStatus)
		private.GET("/stats", GetStats)
		private.GET("/system/log", GetSystemLog)

		// 日志
		private.GET("/log/slinx", SlinxLog)
		private.GET("/log/core", CoreLog)

		// 检测
		private.POST("/detect/ip/fetch", FetchIP)
		private.GET("/detect/ip", DetectIP)
		private.POST("/detect/unlock/fetch", FetchUnlock)
		private.GET("/detect/unlock", DetectUnlock)
		private.GET("/detect/back-route", DetectBackRoute)
		private.POST("/detect/back-route/fetch", FetchBackRoute)

		// WARP
		private.GET("/warp", GetWarp)
		private.POST("/warp/register", RegisterWarp)
		private.POST("/warp/refresh", RefreshWarp)
		private.DELETE("/warp", DeleteWarp)
		private.PUT("/warp/auto-update", SetWarpAutoUpdate)
		private.PUT("/warp/license", SetWarpLicense)

		// 端点
		private.GET("/endpoint", GetEndpoints)
		private.PUT("/endpoint/save", SaveEndpoint)
		private.DELETE("/endpoint/:id", DeleteEndpoint)
		private.PUT("/endpoint/:id/toggle", ToggleEndpoint)
		private.POST("/endpoint/warp", CreateWarpEndpoint)

		// 路由
		private.GET("/rule", GetRule)
		private.PUT("/rule", SaveRule)

		// 生成
		private.GET("/generate/port", GeneratePort)
		private.GET("/generate/reality-target", GenerateRealityTarget)
		private.GET("/generate/reality-keypair", GenerateRealityKeyPair)
		private.GET("/generate/shortids", GenerateShortIDs)
		private.GET("/generate/token", GenerateToken)
		private.GET("/generate/uuid", GenerateUUID)
		private.GET("/generate/password", GeneratePassword)
		private.GET("/generate/wireguard-keypair", GenerateWireguardKeyPair)

		// 推送
		private.GET("/task/:id", TaskLog)

		// 快速
		private.POST("/quick", Quick)
	}
}
