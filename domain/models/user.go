package models

import (
	"database/sql/driver"
	"strings"
	"time"
)

// UserProfile ユーザー設定
type UserProfile struct {
	// 基本情報
	UserID      string  `json:"user_id" gorm:"primaryKey;type:uuid" validate:"required"`
	DisplayName *string `json:"display_name" gorm:"size:100" validate:"omitempty,max=100"`

	// 転職希望条件
	DesiredSalaryMin      *int    `json:"desired_salary_min" gorm:"" validate:"omitempty,min=0"`
	DesiredSalaryMax      *int    `json:"desired_salary_max" gorm:"" validate:"omitempty,min=0"`
	DesiredLocation       *string `json:"desired_location" gorm:"size:200" validate:"omitempty,max=200"`
	DesiredEmploymentType *string `json:"desired_employment_type" gorm:"size:50" validate:"omitempty,max=50"`

	// 検索・フィルタ設定
	DefaultKeywords   *string `json:"default_keywords" gorm:"type:text"`   // カンマ区切りキーワード
	ExcludedCompanies *string `json:"excluded_companies" gorm:"type:text"` // 除外企業リスト
	ActivePlatforms   *string `json:"active_platforms" gorm:"type:text"`   // 利用中プラットフォーム（カンマ区切り）

	// 機能設定
	DuplicateDetectionEnabled bool `json:"duplicate_detection_enabled" gorm:"default:true"`
	EmailNotifications        bool `json:"email_notifications" gorm:"default:true"`
	AutoScrapingEnabled       bool `json:"auto_scraping_enabled" gorm:"default:false"`

	// メタデータ
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// UserProfileRequest ユーザー設定更新リクエスト
type UserProfileRequest struct {
	DisplayName               *string `json:"display_name" validate:"omitempty,max=100"`
	DesiredSalaryMin          *int    `json:"desired_salary_min" validate:"omitempty,min=0"`
	DesiredSalaryMax          *int    `json:"desired_salary_max" validate:"omitempty,min=0"`
	DesiredLocation           *string `json:"desired_location" validate:"omitempty,max=200"`
	DesiredEmploymentType     *string `json:"desired_employment_type" validate:"omitempty,max=50"`
	DefaultKeywords           *string `json:"default_keywords"`
	ExcludedCompanies         *string `json:"excluded_companies"`
	ActivePlatforms           *string `json:"active_platforms"`
	DuplicateDetectionEnabled *bool   `json:"duplicate_detection_enabled"`
	EmailNotifications        *bool   `json:"email_notifications"`
	AutoScrapingEnabled       *bool   `json:"auto_scraping_enabled"`
}

// UserProfileResponse ユーザー設定レスポンス（API用）
type UserProfileResponse struct {
	Profile struct {
		DisplayName               *string  `json:"display_name"`
		DesiredSalaryMin          *int     `json:"desired_salary_min"`
		DesiredSalaryMax          *int     `json:"desired_salary_max"`
		DesiredLocation           *string  `json:"desired_location"`
		DesiredEmploymentType     *string  `json:"desired_employment_type"`
		DefaultKeywords           *string  `json:"default_keywords"`
		ExcludedCompanies         *string  `json:"excluded_companies"`
		ActivePlatforms           []string `json:"active_platforms"`
		DuplicateDetectionEnabled bool     `json:"duplicate_detection_enabled"`
		EmailNotifications        bool     `json:"email_notifications"`
		AutoScrapingEnabled       bool     `json:"auto_scraping_enabled"`
	} `json:"profile"`
}

// ToResponse UserProfile → UserProfileResponse 変換
func (u *UserProfile) ToResponse() UserProfileResponse {
	return UserProfileResponse{
		Profile: struct {
			DisplayName               *string  `json:"display_name"`
			DesiredSalaryMin          *int     `json:"desired_salary_min"`
			DesiredSalaryMax          *int     `json:"desired_salary_max"`
			DesiredLocation           *string  `json:"desired_location"`
			DesiredEmploymentType     *string  `json:"desired_employment_type"`
			DefaultKeywords           *string  `json:"default_keywords"`
			ExcludedCompanies         *string  `json:"excluded_companies"`
			ActivePlatforms           []string `json:"active_platforms"`
			DuplicateDetectionEnabled bool     `json:"duplicate_detection_enabled"`
			EmailNotifications        bool     `json:"email_notifications"`
			AutoScrapingEnabled       bool     `json:"auto_scraping_enabled"`
		}{
			DisplayName:               u.DisplayName,
			DesiredSalaryMin:          u.DesiredSalaryMin,
			DesiredSalaryMax:          u.DesiredSalaryMax,
			DesiredLocation:           u.DesiredLocation,
			DesiredEmploymentType:     u.DesiredEmploymentType,
			DefaultKeywords:           u.DefaultKeywords,
			ExcludedCompanies:         u.ExcludedCompanies,
			ActivePlatforms:           u.GetActivePlatformsList(),
			DuplicateDetectionEnabled: u.DuplicateDetectionEnabled,
			EmailNotifications:        u.EmailNotifications,
			AutoScrapingEnabled:       u.AutoScrapingEnabled,
		},
	}
}

// GetDefaultKeywordsList カンマ区切りキーワードを配列に変換
func (u *UserProfile) GetDefaultKeywordsList() []string {
	if u.DefaultKeywords == nil || *u.DefaultKeywords == "" {
		return []string{}
	}

	keywords := strings.Split(*u.DefaultKeywords, ",")
	result := make([]string, 0, len(keywords))

	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			result = append(result, keyword)
		}
	}

	return result
}

// SetDefaultKeywordsList 配列をカンマ区切り文字列に変換
func (u *UserProfile) SetDefaultKeywordsList(keywords []string) {
	if len(keywords) == 0 {
		u.DefaultKeywords = nil
		return
	}

	// 空文字を除外
	filtered := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			filtered = append(filtered, keyword)
		}
	}

	if len(filtered) == 0 {
		u.DefaultKeywords = nil
		return
	}

	result := strings.Join(filtered, ",")
	u.DefaultKeywords = &result
}

// GetExcludedCompaniesList カンマ区切り除外企業を配列に変換
func (u *UserProfile) GetExcludedCompaniesList() []string {
	if u.ExcludedCompanies == nil || *u.ExcludedCompanies == "" {
		return []string{}
	}

	companies := strings.Split(*u.ExcludedCompanies, ",")
	result := make([]string, 0, len(companies))

	for _, company := range companies {
		company = strings.TrimSpace(company)
		if company != "" {
			result = append(result, company)
		}
	}

	return result
}

// SetExcludedCompaniesList 配列をカンマ区切り文字列に変換
func (u *UserProfile) SetExcludedCompaniesList(companies []string) {
	if len(companies) == 0 {
		u.ExcludedCompanies = nil
		return
	}

	// 空文字を除外
	filtered := make([]string, 0, len(companies))
	for _, company := range companies {
		company = strings.TrimSpace(company)
		if company != "" {
			filtered = append(filtered, company)
		}
	}

	if len(filtered) == 0 {
		u.ExcludedCompanies = nil
		return
	}

	result := strings.Join(filtered, ",")
	u.ExcludedCompanies = &result
}

// GetActivePlatformsList カンマ区切りプラットフォームを配列に変換
func (u *UserProfile) GetActivePlatformsList() []string {
	if u.ActivePlatforms == nil || *u.ActivePlatforms == "" {
		return []string{}
	}

	platforms := strings.Split(*u.ActivePlatforms, ",")
	result := make([]string, 0, len(platforms))

	for _, platform := range platforms {
		platform = strings.TrimSpace(platform)
		if platform != "" {
			result = append(result, platform)
		}
	}

	return result
}

// SetActivePlatformsList 配列をカンマ区切り文字列に変換
func (u *UserProfile) SetActivePlatformsList(platforms []string) {
	if len(platforms) == 0 {
		u.ActivePlatforms = nil
		return
	}

	// 空文字を除外
	filtered := make([]string, 0, len(platforms))
	for _, platform := range platforms {
		platform = strings.TrimSpace(platform)
		if platform != "" {
			filtered = append(filtered, platform)
		}
	}

	if len(filtered) == 0 {
		u.ActivePlatforms = nil
		return
	}

	result := strings.Join(filtered, ",")
	u.ActivePlatforms = &result
}

// Value SQLドライバー用
func (u UserProfile) Value() (driver.Value, error) {
	return u.UserID, nil
}

// User 認証用ユーザーエンティティ
type User struct {
	ID           string     `json:"id" gorm:"primaryKey;type:uuid"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string     `json:"-" gorm:"not null"` // JSONには含めない
	DisplayName  *string    `json:"display_name" gorm:"size:100"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName テーブル名を指定
func (User) TableName() string {
	return "users"
}
