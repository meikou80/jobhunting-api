package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"jobhunting-api/domain/validators"
	"jobhunting-api/infra/database"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// バリデーターの初期化
	validators.Init()

	// データベース接続
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Ginエンジンの初期化
	r := gin.Default()

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		// データベース接続チェック
		dbStatus := "ok"
		if err := database.Health(); err != nil {
			dbStatus = "error"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"service":  "jobhunting-api",
			"version":  "0.1.0",
			"database": dbStatus,
		})
	})

	// CORSミドルウェア（開発用）
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// API v1 グループ
	v1 := r.Group("/api/v1")
	{
		// 今後ここにエンドポイントを追加していく
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server starting on port %s\n", port)
	log.Fatal(r.Run(":" + port))
}
