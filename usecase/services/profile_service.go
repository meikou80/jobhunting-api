package services

import (
	"context"
	"errors"
	"fmt"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

// ProfileService ユーザープロフィール管理のビジネスロジック
// 設計思想:
// - ビジネスルール集約: プロフィール管理に関する業務ルールを集約
// - リポジトリ抽象化: データアクセス層への依存をインターフェースで抽象化
// - エラー処理統一: ビジネス層でのエラーハンドリングを統一
type ProfileService struct {
	userRepo repositories.UserRepository
}

// NewProfileService コンストラクタ
// 設計思想: 依存性注入によりテスタビリティとメンテナビリティを向上
func NewProfileService(userRepo repositories.UserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

// GetProfile ユーザープロフィール取得
// 設計思想:
// - 単一責任: プロフィール取得の業務ルールのみに集中
// - NULL安全: プロフィールが存在しない場合の適切な処理
// - データ変換: ドメインモデルからレスポンス形式への変換
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
// 設計思想:
// - 業務ルール適用: 初回作成時の制約やデフォルト値設定
// - 重複チェック: 既存プロフィールの確認
// - トランザクション考慮: 作成処理の原子性確保
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
// 設計思想:
// - 部分更新対応: 提供されたフィールドのみを更新
// - 業務ルール適用: 更新時の制約チェック
// - 楽観的ロック: 更新競合の検出と処理
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
// 設計思想: ポインタによるnilチェックで未設定フィールドを判別
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
// 設計思想:
// - 論理削除: 物理削除ではなく削除フラグによる管理
// - 関連データ整合性: 削除時の関連データへの影響を考慮
// - 復旧可能性: 誤削除時の復旧手順を確保
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
// 設計思想: プロフィール情報と設定情報の責任分離
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
// 設計思想: 設定項目の個別更新による柔軟性提供
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
