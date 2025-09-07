package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "jobhunting-api/docs" // swag initで生成される
	"jobhunting-api/infra/database"
	"jobhunting-api/infra/router"
)

// @title           転職活動管理API
// @version         1.0
// @description     複数転職サービスの一元管理と重複防止を実現するAPI
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// データベース接続
	if err := database.Init(); err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// DBマイグレーション実行
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// ルーター初期化
	r := router.NewRouter(database.DB)
	r.SetupRoutes()

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(r.Run(":" + port))
}
