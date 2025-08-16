package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"jobhunting-api/domain/models"
)

var DB *gorm.DB

// Config データベース設定
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfig 環境変数からDB設定を読み込み
func LoadConfig() *Config {
	return &Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "jobhunting_db"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// Init データベース接続を初期化
func Init() error {
	config := LoadConfig()

	// データソース名（DSN）を構築
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Tokyo",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)

	// GORM設定
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel()),
		NowFunc: func() time.Time {
			return time.Now().In(time.FixedZone("JST", 9*60*60))
		},
	}

	// データベース接続
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 接続プールの設定
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 接続プール設定
	sqlDB.SetMaxIdleConns(10)           // アイドル接続の最大数
	sqlDB.SetMaxOpenConns(100)          // オープン接続の最大数
	sqlDB.SetConnMaxLifetime(time.Hour) // 接続の最大ライフタイム

	log.Println("✅ Database connection established")
	return nil
}

// AutoMigrate テーブル自動作成（開発用）
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// モデルの自動マイグレーション
	err := DB.AutoMigrate(
		&models.Job{},
		&models.Application{},
		&models.UserProfile{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("✅ Database migration completed")
	return nil
}

// Close データベース接続を閉じる
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Health データベース接続ヘルスチェック
func Health() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// getEnv 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getLogLevel 環境に応じたログレベルを取得
func getLogLevel() logger.LogLevel {
	switch os.Getenv("GIN_MODE") {
	case "release":
		return logger.Error
	case "test":
		return logger.Silent
	default:
		return logger.Info
	}
}
