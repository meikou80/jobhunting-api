package repositories

import (
	"context"

	"jobhunting-api/domain/models"
)

// ApplicationRepository 応募リポジトリインターフェース
type ApplicationRepository interface {
	// 基本CRUD
	Create(ctx context.Context, application *models.Application) error
	GetByID(ctx context.Context, id int) (*models.Application, error)
	GetByIDWithRelations(ctx context.Context, id int) (*models.Application, error)
	Update(ctx context.Context, application *models.Application) error
	Delete(ctx context.Context, id int) error

	// ユーザー関連
	GetByUserID(ctx context.Context, userID string, filter *models.ApplicationFilter) ([]*models.Application, int, error)
	GetByUserIDWithRelations(ctx context.Context, userID string, filter *models.ApplicationFilter) ([]*models.Application, int, error)
	GetByUserAndJobID(ctx context.Context, userID string, jobID int) (*models.Application, error)

	// ステータス管理
	UpdateStatus(ctx context.Context, id int, status models.ApplicationStatus) error
	GetByStatus(ctx context.Context, userID string, status models.ApplicationStatus) ([]*models.Application, error)

	// 統計・集計
	GetSummaryByUser(ctx context.Context, userID string) (map[string]int, error)
	CountByUser(ctx context.Context, userID string) (int64, error)
	CountByUserAndStatus(ctx context.Context, userID string, status models.ApplicationStatus) (int64, error)

	// 期限・アクション管理
	GetUpcomingActions(ctx context.Context, userID string, days int) ([]*models.Application, error)
	GetRecentActivity(ctx context.Context, userID string, limit int) ([]*models.Application, error)
}
