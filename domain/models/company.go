package models

import (
	"database/sql/driver"
	"time"
)

// Company 企業マスタ
type Company struct {
	// 基本情報
	ID           int     `json:"id" db:"id"`
	Name         string  `json:"name" db:"name" validate:"required,max=200"`
	Industry     *string `json:"industry" db:"industry" validate:"omitempty,max=100"`
	SizeCategory *string `json:"size_category" db:"size_category" validate:"omitempty,oneof=startup medium large"`
	Location     *string `json:"location" db:"location" validate:"omitempty,max=200"`
	WebsiteURL   *string `json:"website_url" db:"website_url" validate:"omitempty,url"`

	// 個人メモ
	Notes  *string `json:"notes" db:"notes"`
	Rating *int    `json:"rating" db:"rating" validate:"omitempty,min=1,max=5"`

	// メタデータ
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// CompanyRequest 企業登録・更新リクエスト
type CompanyRequest struct {
	Name         string  `json:"name" validate:"required,max=200"`
	Industry     *string `json:"industry" validate:"omitempty,max=100"`
	SizeCategory *string `json:"size_category" validate:"omitempty,oneof=startup medium large"`
	Location     *string `json:"location" validate:"omitempty,max=200"`
	WebsiteURL   *string `json:"website_url" validate:"omitempty,url"`
	Notes        *string `json:"notes"`
	Rating       *int    `json:"rating" validate:"omitempty,min=1,max=5"`
}

// CompanyResponse 企業レスポンス（API用）
type CompanyResponse struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Industry     *string `json:"industry"`
	SizeCategory *string `json:"size_category"`
	Location     *string `json:"location"`
	WebsiteURL   *string `json:"website_url"`
	Notes        *string `json:"notes"`
	Rating       *int    `json:"rating"`
	CreatedAt    string  `json:"created_at"` // ISO 8601形式の文字列
	UpdatedAt    string  `json:"updated_at"`
}

// CompanyListResponse 企業一覧レスポンス
type CompanyListResponse struct {
	Companies []CompanyResponse `json:"companies"`
	Total     int               `json:"total"`
	Page      int               `json:"page"`
}

// ToResponse Company → CompanyResponse 変換
func (c *Company) ToResponse() CompanyResponse {
	return CompanyResponse{
		ID:           c.ID,
		Name:         c.Name,
		Industry:     c.Industry,
		SizeCategory: c.SizeCategory,
		Location:     c.Location,
		WebsiteURL:   c.WebsiteURL,
		Notes:        c.Notes,
		Rating:       c.Rating,
		CreatedAt:    c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
	}
}

// CompanyFilter 企業検索フィルタ
type CompanyFilter struct {
	Search   *string `form:"search" validate:"omitempty,max=100"`
	Industry *string `form:"industry" validate:"omitempty,max=100"`
	Size     *string `form:"size" validate:"omitempty,oneof=startup medium large"`
	Page     int     `form:"page" validate:"min=1" default:"1"`
	Limit    int     `form:"limit" validate:"min=1,max=100" default:"20"`
}

// Value SQLドライバー用
func (c Company) Value() (driver.Value, error) {
	return c.ID, nil
}
