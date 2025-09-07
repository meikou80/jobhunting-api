package repositories

import (
	"context"

	"gorm.io/gorm"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

// userRepository ユーザーリポジトリの実装
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository ユーザーリポジトリのコンストラクタ
func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

// CreateProfile ユーザープロフィールを作成
func (r *userRepository) CreateProfile(ctx context.Context, profile *models.UserProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// GetProfile ユーザープロフィールを取得
func (r *userRepository) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile ユーザープロフィールを更新
func (r *userRepository) UpdateProfile(ctx context.Context, profile *models.UserProfile) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", profile.UserID).
		Updates(profile).Error
}

// DeleteProfile ユーザープロフィールを削除
func (r *userRepository) DeleteProfile(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.UserProfile{}).Error
}

// UpdateKeywords デフォルトキーワードを更新
func (r *userRepository) UpdateKeywords(ctx context.Context, userID string, keywords []string) error {
	profile := &models.UserProfile{UserID: userID}
	profile.SetDefaultKeywordsList(keywords)

	return r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("user_id = ?", userID).
		Update("default_keywords", profile.DefaultKeywords).Error
}

// GetKeywords デフォルトキーワードを取得
func (r *userRepository) GetKeywords(ctx context.Context, userID string) ([]string, error) {
	profile, err := r.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile.GetDefaultKeywordsList(), nil
}

// UpdateExcludedCompanies 除外企業を更新
func (r *userRepository) UpdateExcludedCompanies(ctx context.Context, userID string, companies []string) error {
	profile := &models.UserProfile{UserID: userID}
	profile.SetExcludedCompaniesList(companies)

	return r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("user_id = ?", userID).
		Update("excluded_companies", profile.ExcludedCompanies).Error
}

// GetExcludedCompanies 除外企業を取得
func (r *userRepository) GetExcludedCompanies(ctx context.Context, userID string) ([]string, error) {
	profile, err := r.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile.GetExcludedCompaniesList(), nil
}

// UpdateNotificationSettings 通知設定を更新
func (r *userRepository) UpdateNotificationSettings(ctx context.Context, userID string, emailNotifications bool) error {
	return r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("user_id = ?", userID).
		Update("email_notifications", emailNotifications).Error
}

// GetNotificationSettings 通知設定を取得
func (r *userRepository) GetNotificationSettings(ctx context.Context, userID string) (bool, error) {
	var profile models.UserProfile
	err := r.db.WithContext(ctx).
		Select("email_notifications").
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		return false, err
	}

	return profile.EmailNotifications, nil
}

// UpdateActivePlatforms 利用中プラットフォームを更新
func (r *userRepository) UpdateActivePlatforms(ctx context.Context, userID string, platforms []string) error {
	profile := &models.UserProfile{UserID: userID}
	profile.SetActivePlatformsList(platforms)

	return r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("user_id = ?", userID).
		Update("active_platforms", profile.ActivePlatforms).Error
}

// GetActivePlatforms 利用中プラットフォームを取得
func (r *userRepository) GetActivePlatforms(ctx context.Context, userID string) ([]string, error) {
	profile, err := r.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return profile.GetActivePlatformsList(), nil
}

// UpdateFeatureSettings 機能設定を更新
func (r *userRepository) UpdateFeatureSettings(ctx context.Context, userID string, duplicateDetection, autoScraping bool) error {
	updates := map[string]interface{}{
		"duplicate_detection_enabled": duplicateDetection,
		"auto_scraping_enabled":       autoScraping,
	}

	return r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("user_id = ?", userID).
		Updates(updates).Error
}

// GetFeatureSettings 機能設定を取得
func (r *userRepository) GetFeatureSettings(ctx context.Context, userID string) (duplicateDetection bool, autoScraping bool, err error) {
	var profile models.UserProfile
	err = r.db.WithContext(ctx).
		Select("duplicate_detection_enabled, auto_scraping_enabled").
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		return false, false, err
	}

	return profile.DuplicateDetectionEnabled, profile.AutoScrapingEnabled, nil
}

// Exists ユーザープロフィールの存在チェック
func (r *userRepository) Exists(ctx context.Context, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.UserProfile{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	return count > 0, err
}

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}
