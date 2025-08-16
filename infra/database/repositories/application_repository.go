package repositories

import (
	"context"

	"gorm.io/gorm"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

// applicationRepository 応募リポジトリの実装
type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository 応募リポジトリのコンストラクタ
func NewApplicationRepository(db *gorm.DB) repositories.ApplicationRepository {
	return &applicationRepository{db: db}
}

// Create 応募を作成
func (r *applicationRepository) Create(ctx context.Context, application *models.Application) error {
	return r.db.WithContext(ctx).Create(application).Error
}

// GetByID IDで応募を取得
func (r *applicationRepository) GetByID(ctx context.Context, id int) (*models.Application, error) {
	var application models.Application
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&application).Error

	if err != nil {
		return nil, err
	}

	// ActivityLogの復元
	if err := r.restoreActivityLog(&application); err != nil {
		return nil, err
	}

	return &application, nil
}

// GetByIDWithRelations IDで応募を取得（関連情報込み）
func (r *applicationRepository) GetByIDWithRelations(ctx context.Context, id int) (*models.Application, error) {
	var application models.Application
	err := r.db.WithContext(ctx).
		Preload("Job").
		Where("applications.id = ? AND applications.deleted_at IS NULL", id).
		First(&application).Error

	if err != nil {
		return nil, err
	}

	// ActivityLogの復元
	if err := r.restoreActivityLog(&application); err != nil {
		return nil, err
	}

	return &application, nil
}

// Update 応募を更新
func (r *applicationRepository) Update(ctx context.Context, application *models.Application) error {
	// ActivityLogをJSONに変換
	if err := r.saveActivityLog(application); err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", application.ID).
		Updates(application).Error
}

// Delete 応募を論理削除
func (r *applicationRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).
		Model(&models.Application{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// GetByUserID ユーザーIDで応募一覧を取得
func (r *applicationRepository) GetByUserID(ctx context.Context, userID string, filter *models.ApplicationFilter) ([]*models.Application, int, error) {
	var applications []*models.Application
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Application{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

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
		Order(r.getOrderBy(filter.Sort)).
		Find(&applications).Error

	if err != nil {
		return nil, 0, err
	}

	// ActivityLogの復元
	for _, app := range applications {
		if err := r.restoreActivityLog(app); err != nil {
			return nil, 0, err
		}
	}

	return applications, int(total), nil
}

// GetByUserIDWithRelations ユーザーIDで応募一覧を取得（関連情報込み）
func (r *applicationRepository) GetByUserIDWithRelations(ctx context.Context, userID string, filter *models.ApplicationFilter) ([]*models.Application, int, error) {
	var applications []*models.Application
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Application{}).
		Joins("LEFT JOIN jobs ON applications.job_id = jobs.id").
		Where("applications.user_id = ? AND applications.deleted_at IS NULL", userID)

	// フィルタ適用
	query = r.applyFilters(query, filter)

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Preload("Job").
		Offset(offset).
		Limit(filter.Limit).
		Order(r.getOrderBy(filter.Sort)).
		Find(&applications).Error

	if err != nil {
		return nil, 0, err
	}

	// ActivityLogの復元
	for _, app := range applications {
		if err := r.restoreActivityLog(app); err != nil {
			return nil, 0, err
		}
	}

	return applications, int(total), nil
}

// GetByUserAndJobID ユーザーIDと求人IDで応募を取得
func (r *applicationRepository) GetByUserAndJobID(ctx context.Context, userID string, jobID int) (*models.Application, error) {
	var application models.Application
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND job_id = ? AND deleted_at IS NULL", userID, jobID).
		First(&application).Error

	if err != nil {
		return nil, err
	}

	// ActivityLogの復元
	if err := r.restoreActivityLog(&application); err != nil {
		return nil, err
	}

	return &application, nil
}

// UpdateStatus ステータスを更新
func (r *applicationRepository) UpdateStatus(ctx context.Context, id int, status models.ApplicationStatus) error {
	return r.db.WithContext(ctx).
		Model(&models.Application{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("status", status).Error
}

// GetByStatus ステータスで応募を取得
func (r *applicationRepository) GetByStatus(ctx context.Context, userID string, status models.ApplicationStatus) ([]*models.Application, error) {
	var applications []*models.Application

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, status).
		Order("updated_at DESC").
		Find(&applications).Error

	if err != nil {
		return nil, err
	}

	// ActivityLogの復元
	for _, app := range applications {
		if err := r.restoreActivityLog(app); err != nil {
			return nil, err
		}
	}

	return applications, nil
}

// GetSummaryByUser ユーザーの応募サマリを取得
func (r *applicationRepository) GetSummaryByUser(ctx context.Context, userID string) (map[string]int, error) {
	var results []struct {
		Status string
		Count  int
	}

	err := r.db.WithContext(ctx).
		Model(&models.Application{}).
		Select("status, COUNT(*) as count").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("status").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]int)
	for _, result := range results {
		summary[result.Status] = result.Count
	}

	return summary, nil
}

// CountByUser ユーザーの応募総数を取得
func (r *applicationRepository) CountByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Application{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// CountByUserAndStatus ユーザーの特定ステータスの応募数を取得
func (r *applicationRepository) CountByUserAndStatus(ctx context.Context, userID string, status models.ApplicationStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Application{}).
		Where("user_id = ? AND status = ? AND deleted_at IS NULL", userID, status).
		Count(&count).Error
	return count, err
}

// GetUpcomingActions 今後のアクションが必要な応募を取得
func (r *applicationRepository) GetUpcomingActions(ctx context.Context, userID string, days int) ([]*models.Application, error) {
	// 実装は簡略化：面接ステータスの応募を返す
	return r.GetByStatus(ctx, userID, models.StatusInterview)
}

// GetRecentActivity 最近の活動を取得
func (r *applicationRepository) GetRecentActivity(ctx context.Context, userID string, limit int) ([]*models.Application, error) {
	var applications []*models.Application

	err := r.db.WithContext(ctx).
		Preload("Job").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("updated_at DESC").
		Limit(limit).
		Find(&applications).Error

	if err != nil {
		return nil, err
	}

	// ActivityLogの復元
	for _, app := range applications {
		if err := r.restoreActivityLog(app); err != nil {
			return nil, err
		}
	}

	return applications, nil
}

// applyFilters フィルタを適用
func (r *applicationRepository) applyFilters(query *gorm.DB, filter *models.ApplicationFilter) *gorm.DB {
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("applications.status = ?", *filter.Status)
	}
	if filter.AppliedVia != nil && *filter.AppliedVia != "" {
		query = query.Where("applications.applied_via = ?", *filter.AppliedVia)
	}
	if filter.CompanyName != nil && *filter.CompanyName != "" {
		query = query.Where("jobs.company_name ILIKE ?", "%"+*filter.CompanyName+"%")
	}

	return query
}

// getOrderBy ソート順を取得
func (r *applicationRepository) getOrderBy(sort *string) string {
	if sort == nil {
		return "applications.updated_at DESC"
	}

	switch *sort {
	case "applied_date_desc":
		return "applications.applied_date DESC"
	case "status_asc":
		return "applications.status ASC, applications.updated_at DESC"
	case "updated_desc":
		return "applications.updated_at DESC"
	default:
		return "applications.updated_at DESC"
	}
}

// saveActivityLog ActivityLogをJSONに変換して保存
func (r *applicationRepository) saveActivityLog(application *models.Application) error {
	if application.ActivityLog == nil {
		return nil
	}

	return application.SetActivityLog(application.ActivityLog)
}

// restoreActivityLog JSON文字列からActivityLogを復元
func (r *applicationRepository) restoreActivityLog(application *models.Application) error {
	log, err := application.GetActivityLog()
	if err != nil {
		return err
	}

	application.ActivityLog = log
	return nil
}
