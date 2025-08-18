package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"jobhunting-api/infra/database/repositories"
	"jobhunting-api/infra/middleware"
	"jobhunting-api/interface/controllers"
	"jobhunting-api/usecase/services"
)

type Router struct {
	engine *gin.Engine
	db     *gorm.DB
}

func NewRouter(db *gorm.DB) *Router {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware())

	return &Router{
		engine: engine,
		db:     db,
	}
}

func (r *Router) SetupRoutes() {
	// ヘルスチェック
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Job Hunting API is running",
		})
	})

	// 依存性注入とコントローラー初期化
	r.setupDependencies()
}

// setupDependencies 依存関係の設定とAPIルート構築
// 設計思想: Clean Architectureの依存関係を外側から内側に注入
func (r *Router) setupDependencies() {
	// Repository層の初期化
	repos := repositories.NewRepositories(r.db)

	// Service層の初期化（Repository注入）
	profileService := services.NewProfileService(repos.User)

	// Controller層の初期化（Service注入）
	profileController := controllers.NewProfileController(profileService)

	// Middleware初期化
	authMiddleware := middleware.NewAuthMiddleware()

	// APIルートグループ設定
	r.setupAPIRoutes(authMiddleware, profileController)
}

// setupAPIRoutes API v1ルートの設定
// 設計思想: バージョニング対応とミドルウェア適用の明確化
func (r *Router) setupAPIRoutes(
	authMiddleware *middleware.AuthMiddleware,
	profileController *controllers.ProfileController,
) {
	// API v1グループ
	v1 := r.engine.Group("/api/v1")

	// 認証必須のプロフィール管理エンドポイント
	// 設計思想: 全てのプロフィール操作で認証を必須とする
	profileGroup := v1.Group("/profile")
	profileGroup.Use(authMiddleware.RequireAuth())
	{
		profileGroup.GET("", profileController.GetProfile)       // GET /api/v1/profile
		profileGroup.POST("", profileController.CreateProfile)   // POST /api/v1/profile
		profileGroup.PUT("", profileController.UpdateProfile)    // PUT /api/v1/profile
		profileGroup.DELETE("", profileController.DeleteProfile) // DELETE /api/v1/profile

		// 設定管理サブリソース
		profileGroup.GET("/settings", profileController.GetUserSettings)    // GET /api/v1/profile/settings
		profileGroup.PUT("/settings", profileController.UpdateUserSettings) // PUT /api/v1/profile/settings
	}

	// TODO: Phase 2で以下を追加
	// - /api/v1/jobs (求人管理)
	// - /api/v1/applications (応募管理)
	// - /api/v1/platforms (プラットフォーム管理)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// Run サーバー起動
func (r *Router) Run(addr string) error {
	log.Printf("Server is running on port %s", addr)
	return r.engine.Run(addr)
}
