package services

import (
	"context"
	"errors"
	"fmt"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

// ProfileService ユーザープロフィール管理のビジネスロジック
type ProfileService struct {
	userRepo repositories.UserRepository
}

// NewProfileService コンストラクタ
func NewProfileService(userRepo repositories.UserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

// GetProfile ユーザープロフィール取得
func (ps *ProfileService) GetProfile(userID string) (*models.UserProfileResponse, error) {
	// 入力検証
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// プロフィール取得
	profile, err := ps.userRepo.GetProfile(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if profile == nil {
		return nil, errors.New("profile not found")
	}

	// レスポンス形式に変換
	response := profile.ToResponse()
	return &response, nil
}

// CreateProfile ユーザープロフィール作成
func (ps *ProfileService) CreateProfile(userID string, req *models.UserProfileRequest) (*models.UserProfileResponse, error) {
	// 入力検証
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if req == nil {
		return nil, errors.New("request data is required")
	}

	// 既存プロフィールチェック
	existingProfile, err := ps.userRepo.GetProfile(context.TODO(), userID)
	if err == nil && existingProfile != nil {
		return nil, errors.New("profile already exists")
	}

	// ドメインモデル作成
	profile := &models.UserProfile{
		UserID:                userID,
		DisplayName:           req.DisplayName,
		DesiredSalaryMin:      req.DesiredSalaryMin,
		DesiredSalaryMax:      req.DesiredSalaryMax,
		DesiredLocation:       req.DesiredLocation,
		DesiredEmploymentType: req.DesiredEmploymentType,
		DefaultKeywords:       req.DefaultKeywords,
		ExcludedCompanies:     req.ExcludedCompanies,
		// bool型フィールドはデフォルト値を直接設定
		DuplicateDetectionEnabled: true,
		EmailNotifications:        true,
		AutoScrapingEnabled:       false,
	}

	// リクエストでbool値が指定されている場合は上書き
	if req.DuplicateDetectionEnabled != nil {
		profile.DuplicateDetectionEnabled = *req.DuplicateDetectionEnabled
	}
	if req.EmailNotifications != nil {
		profile.EmailNotifications = *req.EmailNotifications
	}
	if req.AutoScrapingEnabled != nil {
		profile.AutoScrapingEnabled = *req.AutoScrapingEnabled
	}

	// ActivePlatforms の設定
	if req.ActivePlatforms != nil {
		profile.SetActivePlatformsList([]string{*req.ActivePlatforms})
	}

	// プロフィール作成
	err = ps.userRepo.CreateProfile(context.TODO(), profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	response := profile.ToResponse()
	return &response, nil
}

// UpdateProfile ユーザープロフィール更新
func (ps *ProfileService) UpdateProfile(userID string, req *models.UserProfileRequest) (*models.UserProfileResponse, error) {
	// 入力検証
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if req == nil {
		return nil, errors.New("request data is required")
	}

	// 既存プロフィール取得
	existingProfile, err := ps.userRepo.GetProfile(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing profile: %w", err)
	}
	if existingProfile == nil {
		return nil, errors.New("profile not found")
	}

	// 更新フィールドの適用（部分更新）
	ps.applyUpdateFields(existingProfile, req)

	// プロフィール更新
	err = ps.userRepo.UpdateProfile(context.TODO(), existingProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	response := existingProfile.ToResponse()
	return &response, nil
}

// applyUpdateFields 部分更新フィールドの適用
func (ps *ProfileService) applyUpdateFields(profile *models.UserProfile, req *models.UserProfileRequest) {
	if req.DisplayName != nil {
		profile.DisplayName = req.DisplayName
	}
	if req.DesiredSalaryMin != nil {
		profile.DesiredSalaryMin = req.DesiredSalaryMin
	}
	if req.DesiredSalaryMax != nil {
		profile.DesiredSalaryMax = req.DesiredSalaryMax
	}
	if req.DesiredLocation != nil {
		profile.DesiredLocation = req.DesiredLocation
	}
	if req.DesiredEmploymentType != nil {
		profile.DesiredEmploymentType = req.DesiredEmploymentType
	}
	if req.DefaultKeywords != nil {
		profile.DefaultKeywords = req.DefaultKeywords
	}
	if req.ExcludedCompanies != nil {
		profile.ExcludedCompanies = req.ExcludedCompanies
	}
	if req.DuplicateDetectionEnabled != nil {
		profile.DuplicateDetectionEnabled = *req.DuplicateDetectionEnabled
	}
	if req.EmailNotifications != nil {
		profile.EmailNotifications = *req.EmailNotifications
	}
	if req.AutoScrapingEnabled != nil {
		profile.AutoScrapingEnabled = *req.AutoScrapingEnabled
	}
	if req.ActivePlatforms != nil {
		profile.SetActivePlatformsList([]string{*req.ActivePlatforms})
	}
}

// DeleteProfile ユーザープロフィール削除（論理削除）
func (ps *ProfileService) DeleteProfile(userID string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	// 既存プロフィール確認
	existingProfile, err := ps.userRepo.GetProfile(context.TODO(), userID)
	if err != nil {
		return fmt.Errorf("failed to get profile: %w", err)
	}
	if existingProfile == nil {
		return errors.New("profile not found")
	}

	// 論理削除実行
	err = ps.userRepo.DeleteProfile(context.TODO(), userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// GetUserSettings ユーザー設定取得
func (ps *ProfileService) GetUserSettings(userID string) (map[string]interface{}, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// プロフィールから設定情報を抽出
	profile, err := ps.userRepo.GetProfile(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	settings := map[string]interface{}{
		"active_platforms":            profile.GetActivePlatformsList(),
		"duplicate_detection_enabled": profile.DuplicateDetectionEnabled,
		"email_notifications":         profile.EmailNotifications,
		"auto_scraping_enabled":       profile.AutoScrapingEnabled,
	}

	return settings, nil
}

// UpdateUserSettings ユーザー設定更新
func (ps *ProfileService) UpdateUserSettings(userID string, settings map[string]interface{}) (map[string]interface{}, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// 既存プロフィール取得
	profile, err := ps.userRepo.GetProfile(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, errors.New("profile not found")
	}

	// 設定項目の更新
	if activePlatforms, exists := settings["active_platforms"]; exists {
		if platforms, ok := activePlatforms.([]string); ok {
			profile.SetActivePlatformsList(platforms)
		}
	}
	if duplicateDetection, exists := settings["duplicate_detection_enabled"]; exists {
		if enabled, ok := duplicateDetection.(bool); ok {
			profile.DuplicateDetectionEnabled = enabled
		}
	}
	if emailNotifications, exists := settings["email_notifications"]; exists {
		if enabled, ok := emailNotifications.(bool); ok {
			profile.EmailNotifications = enabled
		}
	}
	if autoScraping, exists := settings["auto_scraping_enabled"]; exists {
		if enabled, ok := autoScraping.(bool); ok {
			profile.AutoScrapingEnabled = enabled
		}
	}

	// 更新実行
	err = ps.userRepo.UpdateProfile(context.TODO(), profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	// 更新後の設定を返却
	return ps.GetUserSettings(userID)
}
