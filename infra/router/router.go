package router

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"jobhunting-api/infra/database/repositories"
	"jobhunting-api/infra/middleware"
	"jobhunting-api/interface/controllers"
	"jobhunting-api/usecase/services"

	_ "jobhunting-api/docs" // swag initで生成される
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

	// Swagger UI エンドポイント
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// コントローラー初期化
	r.setupDependencies()
}

// setupDependencies 依存関係の設定とAPIルート構築
func (r *Router) setupDependencies() {
	// Repository層の初期化
	repos := repositories.NewRepositories(r.db)

	// Service層の初期化（Repository注入）
	profileService := services.NewProfileService(repos.User)
	jobService := services.NewJobService(repos.Job)
	authService := services.NewAuthService(repos.Auth) // AuthServiceを追加

	// Controller層の初期化（Service注入）
	profileController := controllers.NewProfileController(profileService)
	jobController := controllers.NewJobController(jobService)
	authController := controllers.NewAuthController(authService) // AuthControllerを追加

	// Middleware初期化
	authMiddleware := middleware.NewAuthMiddleware()

	// APIルートグループ設定
	r.setupAPIRoutes(authMiddleware, profileController, jobController, authController)
}

// setupAPIRoutes API v1ルートの設定
func (r *Router) setupAPIRoutes(
	authMiddleware *middleware.AuthMiddleware,
	profileController *controllers.ProfileController,
	jobController *controllers.JobController,
	authController *controllers.AuthController,
) {
	// API v1グループ
	v1 := r.engine.Group("/api/v1")

	// 認証必須のプロフィール管理エンドポイント
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

	// Phase 2: 求人管理エンドポイント
	jobGroup := v1.Group("/jobs")
	jobGroup.Use(authMiddleware.RequireAuth())
	{
		// 基本CRUD操作
		jobGroup.POST("", jobController.CreateJob)       // POST /api/v1/jobs - 求人登録（重複チェック付き）
		jobGroup.GET("", jobController.GetJobs)          // GET /api/v1/jobs - 求人一覧取得（フィルタ対応）
		jobGroup.GET("/:id", jobController.GetJobDetail) // GET /api/v1/jobs/:id - 求人詳細取得
		jobGroup.PUT("/:id", jobController.UpdateJob)    // PUT /api/v1/jobs/:id - 求人更新
		jobGroup.DELETE("/:id", jobController.DeleteJob) // DELETE /api/v1/jobs/:id - 求人削除

		// Phase 2の核心機能
		jobGroup.POST("/check-duplicates", jobController.CheckDuplicates) // POST /api/v1/jobs/check-duplicates - 重複チェック
		jobGroup.GET("/stats", jobController.GetPlatformStats)            // GET /api/v1/jobs/stats - プラットフォーム統計
		jobGroup.GET("/dashboard", jobController.GetDashboard)            // GET /api/v1/jobs/dashboard - ダッシュボード
	}

	// TODO: Phase 3で以下を追加
	// - /api/v1/applications (応募管理)

	// 認証エンドポイント（認証不要）
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", authController.Register) // POST /api/v1/auth/register
		authGroup.POST("/login", authController.Login)       // POST /api/v1/auth/login
		authGroup.POST("/logout", authController.Logout)     // POST /api/v1/auth/logout
	}

	// 認証必須の認証エンドポイント
	authSecureGroup := v1.Group("/auth")
	authSecureGroup.Use(authMiddleware.RequireAuth())
	{
		authSecureGroup.GET("/me", authController.GetMe)              // GET /api/v1/auth/me
		authSecureGroup.POST("/refresh", authController.RefreshToken) // POST /api/v1/auth/refresh
	}
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
