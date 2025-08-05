package repositories

import (
	"context"

	"jobhunting-api/domain/models"
)

// JobRepository 求人リポジトリインターフェース
type JobRepository interface {
	// 基本CRUD
	Create(ctx context.Context, job *models.Job) error
	GetByID(ctx context.Context, id int) (*models.Job, error)
	GetByIDWithCompany(ctx context.Context, id int) (*models.Job, error)
	Update(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id int) error

	// 検索・フィルタ
	List(ctx context.Context, filter *models.JobFilter) ([]*models.Job, int, error)
	ListWithCompany(ctx context.Context, filter *models.JobFilter) ([]*models.Job, int, error)
	Search(ctx context.Context, query string) ([]*models.Job, error)
	GetByCompanyID(ctx context.Context, companyID int) ([]*models.Job, error)

	// 外部連携
	GetByExternalID(ctx context.Context, sourceSite, externalID string) (*models.Job, error)
	CreateOrUpdate(ctx context.Context, job *models.Job) error

	// 集計
	Count(ctx context.Context) (int64, error)
	CountByCompany(ctx context.Context, companyID int) (int64, error)
}
