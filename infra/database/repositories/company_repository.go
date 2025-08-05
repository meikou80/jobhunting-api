package repositories

import (
	"context"

	"gorm.io/gorm"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

// companyRepository 企業リポジトリの実装
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository 企業リポジトリのコンストラクタ
func NewCompanyRepository(db *gorm.DB) repositories.CompanyRepository {
	return &companyRepository{db: db}
}

// Create 企業を作成
func (r *companyRepository) Create(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).Create(company).Error
}

// GetByID IDで企業を取得
func (r *companyRepository) GetByID(ctx context.Context, id int) (*models.Company, error) {
	var company models.Company
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&company).Error

	if err != nil {
		return nil, err
	}
	return &company, nil
}

// Update 企業を更新
func (r *companyRepository) Update(ctx context.Context, company *models.Company) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", company.ID).
		Updates(company).Error
}

// Delete 企業を論理削除
func (r *companyRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).
		Model(&models.Company{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// List 企業一覧を取得（フィルタ付き）
func (r *companyRepository) List(ctx context.Context, filter *models.CompanyFilter) ([]*models.Company, int, error) {
	var companies []*models.Company
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Company{}).Where("deleted_at IS NULL")

	// フィルタ適用
	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+*filter.Search+"%")
	}
	if filter.Industry != nil && *filter.Industry != "" {
		query = query.Where("industry = ?", *filter.Industry)
	}
	if filter.Size != nil && *filter.Size != "" {
		query = query.Where("size_category = ?", *filter.Size)
	}

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Offset(offset).
		Limit(filter.Limit).
		Order("created_at DESC").
		Find(&companies).Error

	return companies, int(total), err
}

// Search 企業名で検索
func (r *companyRepository) Search(ctx context.Context, query string) ([]*models.Company, error) {
	var companies []*models.Company

	err := r.db.WithContext(ctx).
		Where("name ILIKE ? AND deleted_at IS NULL", "%"+query+"%").
		Order("name ASC").
		Limit(50).
		Find(&companies).Error

	return companies, err
}

// GetByName 企業名で取得
func (r *companyRepository) GetByName(ctx context.Context, name string) (*models.Company, error) {
	var company models.Company
	err := r.db.WithContext(ctx).
		Where("name = ? AND deleted_at IS NULL", name).
		First(&company).Error

	if err != nil {
		return nil, err
	}
	return &company, nil
}

// Count 企業総数を取得
func (r *companyRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Company{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}
