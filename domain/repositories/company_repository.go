package repositories

import (
	"context"

	"jobhunting-api/domain/models"
)

// CompanyRepository 企業リポジトリインターフェース
type CompanyRepository interface {
	// 基本CRUD
	Create(ctx context.Context, company *models.Company) error
	GetByID(ctx context.Context, id int) (*models.Company, error)
	Update(ctx context.Context, company *models.Company) error
	Delete(ctx context.Context, id int) error

	// 検索・フィルタ
	List(ctx context.Context, filter *models.CompanyFilter) ([]*models.Company, int, error)
	Search(ctx context.Context, query string) ([]*models.Company, error)
	GetByName(ctx context.Context, name string) (*models.Company, error)

	// 集計
	Count(ctx context.Context) (int64, error)
}
