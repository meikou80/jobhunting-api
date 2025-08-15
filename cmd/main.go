package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"jobhunting-api/infra/database"
	"jobhunting-api/infra/router"
)

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
