package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/novablog/server/internal/config"
	"github.com/novablog/server/internal/database"
	"github.com/novablog/server/internal/handlers"
	"github.com/novablog/server/internal/middleware"
	"github.com/novablog/server/internal/utils"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置运行模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 创建 JWT 管理器
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	// 创建处理器
	authHandler := handlers.NewAuthHandler(jwtManager)
	commentHandler := handlers.NewCommentHandler()
	likeHandler := handlers.NewLikeHandler()
	microHandler := handlers.NewMicroHandler()

	// 创建路由
	r := gin.New()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS 配置
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "NovaBlog API is running",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 公开接口
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth/me", authHandler.GetCurrentUser) // 获取当前用户信息（需要 token 但通过 header 传递）
		api.GET("/comments", commentHandler.GetComments)
		api.GET("/likes", likeHandler.GetLikeStatus)
		api.POST("/likes", likeHandler.ToggleLike) // 允许访客点赞（基于 IP Hash）

		// 微语公开接口
		api.GET("/micros", microHandler.GetMicros)
		api.GET("/micros/stats", microHandler.GetStats)
		api.GET("/micros/heatmap", microHandler.GetHeatmap)
		api.GET("/micros/:id", microHandler.GetMicro)

		// 需要认证的接口
		authGroup := api.Group("")
		authGroup.Use(middleware.AuthMiddleware(jwtManager))
		{
			// 用户相关
			authGroup.GET("/auth/profile", authHandler.GetProfile)
			authGroup.PUT("/auth/profile", authHandler.UpdateProfile)

			// 评论相关（需要登录才能评论）
			authGroup.POST("/comments", commentHandler.CreateComment)
			authGroup.DELETE("/comments/:id", commentHandler.DeleteComment)

			// 微语相关（需要登录）
			authGroup.POST("/micros", microHandler.CreateMicro)
			authGroup.PUT("/micros/:id", microHandler.UpdateMicro)
			authGroup.DELETE("/micros/:id", microHandler.DeleteMicro)
			authGroup.POST("/micros/:id/like", microHandler.ToggleLike)
		}

		// 管理员接口
		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AuthMiddleware(jwtManager))
		adminGroup.Use(middleware.AdminMiddleware())
		{
			// 管理员接口（未来扩展）
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}