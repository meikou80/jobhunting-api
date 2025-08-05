package repositories

import (
	"context"

	"gorm.io/gorm"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

// jobRepository 求人リポジトリの実装
type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository 求人リポジトリのコンストラクタ
func NewJobRepository(db *gorm.DB) repositories.JobRepository {
	return &jobRepository{db: db}
}

// Create 求人を作成
func (r *jobRepository) Create(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

// GetByID IDで求人を取得
func (r *jobRepository) GetByID(ctx context.Context, id int) (*models.Job, error) {
	var job models.Job
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&job).Error

	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetByIDWithCompany IDで求人を取得（企業情報込み）
func (r *jobRepository) GetByIDWithCompany(ctx context.Context, id int) (*models.Job, error) {
	var job models.Job
	err := r.db.WithContext(ctx).
		Preload("Company").
		Where("jobs.id = ? AND jobs.deleted_at IS NULL", id).
		First(&job).Error

	if err != nil {
		return nil, err
	}
	return &job, nil
}

// Update 求人を更新
func (r *jobRepository) Update(ctx context.Context, job *models.Job) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", job.ID).
		Updates(job).Error
}

// Delete 求人を論理削除
func (r *jobRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).
		Model(&models.Job{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// List 求人一覧を取得（フィルタ付き）
func (r *jobRepository) List(ctx context.Context, filter *models.JobFilter) ([]*models.Job, int, error) {
	var jobs []*models.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Job{}).Where("deleted_at IS NULL")

	// フィルタ適用
	query = r.applyFilters(query, filter)

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Offset(offset).
		Limit(filter.Limit).
		Order("posted_date DESC, created_at DESC").
		Find(&jobs).Error

	return jobs, int(total), err
}

// ListWithCompany 求人一覧を取得（企業情報込み）
func (r *jobRepository) ListWithCompany(ctx context.Context, filter *models.JobFilter) ([]*models.Job, int, error) {
	var jobs []*models.Job
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Job{}).
		Joins("LEFT JOIN companies ON jobs.company_id = companies.id").
		Where("jobs.deleted_at IS NULL AND (companies.deleted_at IS NULL OR companies.id IS NULL)")

	// フィルタ適用
	query = r.applyFilters(query, filter)

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Preload("Company").
		Offset(offset).
		Limit(filter.Limit).
		Order("jobs.posted_date DESC, jobs.created_at DESC").
		Find(&jobs).Error

	return jobs, int(total), err
}

// Search 求人をキーワード検索
func (r *jobRepository) Search(ctx context.Context, query string) ([]*models.Job, error) {
	var jobs []*models.Job

	err := r.db.WithContext(ctx).
		Where("(title ILIKE ? OR description ILIKE ?) AND deleted_at IS NULL", "%"+query+"%", "%"+query+"%").
		Order("posted_date DESC").
		Limit(50).
		Find(&jobs).Error

	return jobs, err
}

// GetByCompanyID 企業IDで求人を取得
func (r *jobRepository) GetByCompanyID(ctx context.Context, companyID int) ([]*models.Job, error) {
	var jobs []*models.Job

	err := r.db.WithContext(ctx).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Order("posted_date DESC").
		Find(&jobs).Error

	return jobs, err
}

// GetByExternalID 外部IDで求人を取得
func (r *jobRepository) GetByExternalID(ctx context.Context, sourceSite, externalID string) (*models.Job, error) {
	var job models.Job
	err := r.db.WithContext(ctx).
		Where("source_site = ? AND external_id = ? AND deleted_at IS NULL", sourceSite, externalID).
		First(&job).Error

	if err != nil {
		return nil, err
	}
	return &job, nil
}

// CreateOrUpdate 求人を作成または更新（外部連携用）
func (r *jobRepository) CreateOrUpdate(ctx context.Context, job *models.Job) error {
	if job.SourceSite != nil && job.ExternalID != nil {
		// 既存の求人をチェック
		existing, err := r.GetByExternalID(ctx, *job.SourceSite, *job.ExternalID)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existing != nil {
			// 更新
			job.ID = existing.ID
			return r.Update(ctx, job)
		}
	}

	// 新規作成
	return r.Create(ctx, job)
}

// Count 求人総数を取得
func (r *jobRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Job{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

// CountByCompany 企業別求人数を取得
func (r *jobRepository) CountByCompany(ctx context.Context, companyID int) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Job{}).
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Count(&count).Error
	return count, err
}

// applyFilters フィルタを適用
func (r *jobRepository) applyFilters(query *gorm.DB, filter *models.JobFilter) *gorm.DB {
	if filter.Keyword != nil && *filter.Keyword != "" {
		query = query.Where("(title ILIKE ? OR description ILIKE ?)", "%"+*filter.Keyword+"%", "%"+*filter.Keyword+"%")
	}
	if filter.CompanyID != nil {
		query = query.Where("company_id = ?", *filter.CompanyID)
	}
	if filter.Location != nil && *filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+*filter.Location+"%")
	}
	if filter.SalaryMin != nil {
		query = query.Where("salary_max IS NULL OR salary_max >= ?", *filter.SalaryMin)
	}
	if filter.EmploymentType != nil && *filter.EmploymentType != "" {
		query = query.Where("employment_type = ?", *filter.EmploymentType)
	}
	if filter.RemoteOption != nil && *filter.RemoteOption != "" {
		query = query.Where("remote_option = ?", *filter.RemoteOption)
	}

	return query
}
