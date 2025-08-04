package models

import (
	"database/sql/driver"
	"strings"
	"time"
)

// UserProfile ユーザー設定
type UserProfile struct {
	// 基本情報
	UserID      string  `json:"user_id" db:"user_id"`
	DisplayName *string `json:"display_name" db:"display_name" validate:"omitempty,max=100"`

	// 転職希望条件
	DesiredSalaryMin      *int    `json:"desired_salary_min" db:"desired_salary_min" validate:"omitempty,min=0"`
	DesiredSalaryMax      *int    `json:"desired_salary_max" db:"desired_salary_max" validate:"omitempty,min=0"`
	DesiredLocation       *string `json:"desired_location" db:"desired_location" validate:"omitempty,max=200"`
	DesiredEmploymentType *string `json:"desired_employment_type" db:"desired_employment_type" validate:"omitempty,max=50"`

	// 検索条件
	DefaultKeywords   *string `json:"default_keywords" db:"default_keywords"`     // カンマ区切り
	ExcludedCompanies *string `json:"excluded_companies" db:"excluded_companies"` // カンマ区切り

	// 通知設定
	EmailNotifications bool `json:"email_notifications" db:"email_notifications"`

	// メタデータ
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserProfileRequest ユーザー設定更新リクエスト
type UserProfileRequest struct {
	DisplayName           *string `json:"display_name" validate:"omitempty,max=100"`
	DesiredSalaryMin      *int    `json:"desired_salary_min" validate:"omitempty,min=0"`
	DesiredSalaryMax      *int    `json:"desired_salary_max" validate:"omitempty,min=0"`
	DesiredLocation       *string `json:"desired_location" validate:"omitempty,max=200"`
	DesiredEmploymentType *string `json:"desired_employment_type" validate:"omitempty,max=50"`
	DefaultKeywords       *string `json:"default_keywords"`
	ExcludedCompanies     *string `json:"excluded_companies"`
	EmailNotifications    *bool   `json:"email_notifications"`
}

// UserProfileResponse ユーザー設定レスポンス（API用）
type UserProfileResponse struct {
	Profile struct {
		DisplayName           *string `json:"display_name"`
		DesiredSalaryMin      *int    `json:"desired_salary_min"`
		DesiredSalaryMax      *int    `json:"desired_salary_max"`
		DesiredLocation       *string `json:"desired_location"`
		DesiredEmploymentType *string `json:"desired_employment_type"`
		DefaultKeywords       *string `json:"default_keywords"`
		ExcludedCompanies     *string `json:"excluded_companies"`
		EmailNotifications    bool    `json:"email_notifications"`
	} `json:"profile"`
}

// ToResponse UserProfile → UserProfileResponse 変換
func (u *UserProfile) ToResponse() UserProfileResponse {
	return UserProfileResponse{
		Profile: struct {
			DisplayName           *string `json:"display_name"`
			DesiredSalaryMin      *int    `json:"desired_salary_min"`
			DesiredSalaryMax      *int    `json:"desired_salary_max"`
			DesiredLocation       *string `json:"desired_location"`
			DesiredEmploymentType *string `json:"desired_employment_type"`
			DefaultKeywords       *string `json:"default_keywords"`
			ExcludedCompanies     *string `json:"excluded_companies"`
			EmailNotifications    bool    `json:"email_notifications"`
		}{
			DisplayName:           u.DisplayName,
			DesiredSalaryMin:      u.DesiredSalaryMin,
			DesiredSalaryMax:      u.DesiredSalaryMax,
			DesiredLocation:       u.DesiredLocation,
			DesiredEmploymentType: u.DesiredEmploymentType,
			DefaultKeywords:       u.DefaultKeywords,
			ExcludedCompanies:     u.ExcludedCompanies,
			EmailNotifications:    u.EmailNotifications,
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

// Value SQLドライバー用
func (u UserProfile) Value() (driver.Value, error) {
	return u.UserID, nil
}
