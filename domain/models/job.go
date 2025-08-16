package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Job 求人情報（統合モデル）
type Job struct {
	// 基本情報
	ID            int     `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyName   string  `json:"company_name" gorm:"not null;size:200" validate:"required,max=200"`
	PositionTitle string  `json:"position_title" gorm:"not null;size:300" validate:"required,max=300"`
	Description   *string `json:"description" gorm:"type:text"`
	Requirements  *string `json:"requirements" gorm:"type:text"`

	// プラットフォーム情報（重要）
	SourcePlatform string  `json:"source_platform" gorm:"not null;size:50" validate:"required,max=50"`
	ExternalID     *string `json:"external_id" gorm:"size:200;uniqueIndex:idx_platform_external"`
	SourceURL      *string `json:"source_url" gorm:"type:text" validate:"omitempty,url"`

	// 勤務条件
	SalaryMin      *int    `json:"salary_min" gorm:"index"`
	SalaryMax      *int    `json:"salary_max" gorm:"index"`
	Location       *string `json:"location" gorm:"size:200" validate:"omitempty,max=200"`
	EmploymentType *string `json:"employment_type" gorm:"size:50" validate:"omitempty,oneof=正社員 契約社員 業務委託"`
	RemoteOption   *string `json:"remote_option" gorm:"size:50" validate:"omitempty,oneof=リモート可 ハイブリッド 出社必須"`

	// 日程
	PostedDate   *time.Time `json:"posted_date" gorm:"index"`
	DeadlineDate *time.Time `json:"deadline_date"`
	IsActive     bool       `json:"is_active" gorm:"default:true;index"`

	// 個人管理
	Status        string  `json:"status" gorm:"default:interested;size:50" validate:"omitempty,oneof=interested applied interview offer rejected withdrawn"`
	Priority      int     `json:"priority" gorm:"default:3" validate:"min=1,max=5"`
	PersonalNotes *string `json:"personal_notes" gorm:"type:text"`

	// 重複管理
	DuplicateGroupID    *string  `json:"duplicate_group_id" gorm:"type:uuid;index"`
	IsPrimary           bool     `json:"is_primary" gorm:"default:true"`
	DuplicateConfidence *float64 `json:"duplicate_confidence"`

	// メタデータ
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// JobRequest 求人登録・更新リクエスト
type JobRequest struct {
	CompanyName    string  `json:"company_name" validate:"required,max=200"`
	PositionTitle  string  `json:"position_title" validate:"required,max=300"`
	Description    *string `json:"description"`
	Requirements   *string `json:"requirements"`
	SourcePlatform string  `json:"source_platform" validate:"required,max=50"`
	ExternalID     *string `json:"external_id" validate:"omitempty,max=200"`
	SourceURL      *string `json:"source_url" validate:"omitempty,url"`
	SalaryMin      *int    `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax      *int    `json:"salary_max" validate:"omitempty,min=0"`
	Location       *string `json:"location" validate:"omitempty,max=200"`
	EmploymentType *string `json:"employment_type" validate:"omitempty,oneof=正社員 契約社員 業務委託"`
	RemoteOption   *string `json:"remote_option" validate:"omitempty,oneof=リモート可 ハイブリッド 出社必須"`
	PostedDate     *string `json:"posted_date" validate:"omitempty"`
	DeadlineDate   *string `json:"deadline_date" validate:"omitempty"`
	Priority       *int    `json:"priority" validate:"omitempty,min=1,max=5"`
	PersonalNotes  *string `json:"personal_notes"`
}

// JobResponse 求人レスポンス（API用）
type JobResponse struct {
	ID             int     `json:"id"`
	CompanyName    string  `json:"company_name"`
	PositionTitle  string  `json:"position_title"`
	SourcePlatform string  `json:"source_platform"`
	SalaryRange    *string `json:"salary_range"`
	Location       *string `json:"location"`
	EmploymentType *string `json:"employment_type"`
	PostedDate     *string `json:"posted_date"`
	Status         string  `json:"status"`
	Priority       int     `json:"priority"`
	HasDuplicates  bool    `json:"has_duplicates"`
	DuplicateCount int     `json:"duplicate_count"`
}

// JobDetailResponse 求人詳細レスポンス
type JobDetailResponse struct {
	Job struct {
		ID             int                  `json:"id"`
		CompanyName    string               `json:"company_name"`
		PositionTitle  string               `json:"position_title"`
		Description    *string              `json:"description"`
		Requirements   *string              `json:"requirements"`
		SourcePlatform string               `json:"source_platform"`
		ExternalID     *string              `json:"external_id"`
		SourceURL      *string              `json:"source_url"`
		SalaryMin      *int                 `json:"salary_min"`
		SalaryMax      *int                 `json:"salary_max"`
		Location       *string              `json:"location"`
		EmploymentType *string              `json:"employment_type"`
		RemoteOption   *string              `json:"remote_option"`
		PostedDate     *string              `json:"posted_date"`
		DeadlineDate   *string              `json:"deadline_date"`
		Status         string               `json:"status"`
		Priority       int                  `json:"priority"`
		PersonalNotes  *string              `json:"personal_notes"`
		DuplicateJobs  []DuplicateJobInfo   `json:"duplicate_jobs,omitempty"`
		Application    *ApplicationResponse `json:"application,omitempty"`
	} `json:"job"`
}

// DuplicateJobInfo 重複求人情報
type DuplicateJobInfo struct {
	ID             int     `json:"id"`
	SourcePlatform string  `json:"source_platform"`
	ExternalID     *string `json:"external_id"`
	SourceURL      *string `json:"source_url"`
	Confidence     float64 `json:"confidence"`
}

// JobListResponse 求人一覧レスポンス
type JobListResponse struct {
	Jobs  []JobResponse `json:"jobs"`
	Total int           `json:"total"`
	Page  int           `json:"page"`
}

// ToResponse Job → JobResponse 変換
func (j *Job) ToResponse() JobResponse {
	resp := JobResponse{
		ID:             j.ID,
		CompanyName:    j.CompanyName,
		PositionTitle:  j.PositionTitle,
		SourcePlatform: j.SourcePlatform,
		Location:       j.Location,
		EmploymentType: j.EmploymentType,
		Status:         j.Status,
		Priority:       j.Priority,
		HasDuplicates:  j.DuplicateGroupID != nil,
		DuplicateCount: 0, // 実際の重複数は別途計算が必要
	}

	// 年収範囲
	if j.SalaryMin != nil && j.SalaryMax != nil {
		salaryRange := formatSalaryRange(*j.SalaryMin, *j.SalaryMax)
		resp.SalaryRange = &salaryRange
	}

	// 投稿日
	if j.PostedDate != nil {
		postedDate := j.PostedDate.Format("2006-01-02")
		resp.PostedDate = &postedDate
	}

	return resp
}

// JobFilter 求人検索フィルタ
type JobFilter struct {
	Keyword        *string `form:"keyword" validate:"omitempty,max=100"`
	Platform       *string `form:"platform" validate:"omitempty,max=50"`
	CompanyName    *string `form:"company_name" validate:"omitempty,max=100"`
	Location       *string `form:"location" validate:"omitempty,max=100"`
	SalaryMin      *int    `form:"salary_min" validate:"omitempty,min=0"`
	EmploymentType *string `form:"employment_type" validate:"omitempty,oneof=正社員 契約社員 業務委託"`
	RemoteOption   *string `form:"remote_option" validate:"omitempty,oneof=リモート可 ハイブリッド 出社必須"`
	Status         *string `form:"status" validate:"omitempty,oneof=interested applied interview offer rejected withdrawn"`
	ShowDuplicates *bool   `form:"show_duplicates"`
	Page           int     `form:"page" validate:"min=1" default:"1"`
	Limit          int     `form:"limit" validate:"min=1,max=100" default:"20"`
}

// formatSalaryRange 年収範囲をフォーマット
func formatSalaryRange(min, max int) string {
	return fmt.Sprintf("%d万円〜%d万円", min/10000, max/10000)
}

// Value SQLドライバー用
func (j Job) Value() (driver.Value, error) {
	return j.ID, nil
}
