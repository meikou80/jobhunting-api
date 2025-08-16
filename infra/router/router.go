package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Job Hunting API is running",
		})
	})

	// TODO: Phase 2でAPI routes (/api/v1) を追加
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
