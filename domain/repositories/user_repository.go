package repositories

import (
	"context"

	"jobhunting-api/domain/models"
)

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
	// 基本CRUD
	CreateProfile(ctx context.Context, profile *models.UserProfile) error
	GetProfile(ctx context.Context, userID string) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *models.UserProfile) error
	DeleteProfile(ctx context.Context, userID string) error

	// 設定管理
	UpdateKeywords(ctx context.Context, userID string, keywords []string) error
	GetKeywords(ctx context.Context, userID string) ([]string, error)
	UpdateExcludedCompanies(ctx context.Context, userID string, companies []string) error
	GetExcludedCompanies(ctx context.Context, userID string) ([]string, error)

	// 通知設定
	UpdateNotificationSettings(ctx context.Context, userID string, emailNotifications bool) error
	GetNotificationSettings(ctx context.Context, userID string) (bool, error)

	// 存在チェック
	Exists(ctx context.Context, userID string) (bool, error)
}
