package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Job 求人情報
type Job struct {
	// 基本情報
	ID          int     `json:"id" db:"id"`
	CompanyID   *int    `json:"company_id" db:"company_id"`
	Title       string  `json:"title" db:"title" validate:"required,max=300"`
	Description *string `json:"description" db:"description"`

	// 勤務条件
	SalaryMin      *int    `json:"salary_min" db:"salary_min"`
	SalaryMax      *int    `json:"salary_max" db:"salary_max"`
	Location       *string `json:"location" db:"location" validate:"omitempty,max=200"`
	EmploymentType *string `json:"employment_type" db:"employment_type" validate:"omitempty,oneof=正社員 契約社員 業務委託"`
	RemoteOption   *string `json:"remote_option" db:"remote_option" validate:"omitempty,oneof=リモート可 ハイブリッド 出社必須"`

	// 外部情報
	SourceSite *string `json:"source_site" db:"source_site" validate:"omitempty,max=100"`
	SourceURL  *string `json:"source_url" db:"source_url" validate:"omitempty,url"`
	ExternalID *string `json:"external_id" db:"external_id" validate:"omitempty,max=200"`

	// 日程
	PostedDate   *time.Time `json:"posted_date" db:"posted_date"`
	DeadlineDate *time.Time `json:"deadline_date" db:"deadline_date"`
	IsActive     bool       `json:"is_active" db:"is_active"`

	// リレーション
	Company *Company `json:"company,omitempty"`

	// メタデータ
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// JobRequest 求人登録・更新リクエスト
type JobRequest struct {
	CompanyID      *int    `json:"company_id"`
	Title          string  `json:"title" validate:"required,max=300"`
	Description    *string `json:"description"`
	SalaryMin      *int    `json:"salary_min" validate:"omitempty,min=0"`
	SalaryMax      *int    `json:"salary_max" validate:"omitempty,min=0"`
	Location       *string `json:"location" validate:"omitempty,max=200"`
	EmploymentType *string `json:"employment_type" validate:"omitempty,oneof=正社員 契約社員 業務委託"`
	RemoteOption   *string `json:"remote_option" validate:"omitempty,oneof=リモート可 ハイブリッド 出社必須"`
	SourceSite     *string `json:"source_site" validate:"omitempty,max=100"`
	SourceURL      *string `json:"source_url" validate:"omitempty,url"`
	ExternalID     *string `json:"external_id" validate:"omitempty,max=200"`
	PostedDate     *string `json:"posted_date" validate:"omitempty"`
	DeadlineDate   *string `json:"deadline_date" validate:"omitempty"`
}

// JobResponse 求人レスポンス（API用）
type JobResponse struct {
	ID                int              `json:"id"`
	Title             string           `json:"title"`
	Company           *CompanyResponse `json:"company,omitempty"`
	SalaryRange       *string          `json:"salary_range"`
	Location          *string          `json:"location"`
	EmploymentType    *string          `json:"employment_type"`
	PostedDate        *string          `json:"posted_date"`
	ApplicationStatus *string          `json:"application_status"` // not_applied, applied
	MyPriority        *int             `json:"my_priority"`
}

// JobDetailResponse 求人詳細レスポンス
type JobDetailResponse struct {
	Job struct {
		ID             int                  `json:"id"`
		Title          string               `json:"title"`
		Description    *string              `json:"description"`
		Company        *CompanyResponse     `json:"company"`
		SalaryMin      *int                 `json:"salary_min"`
		SalaryMax      *int                 `json:"salary_max"`
		Location       *string              `json:"location"`
		EmploymentType *string              `json:"employment_type"`
		RemoteOption   *string              `json:"remote_option"`
		SourceSite     *string              `json:"source_site"`
		SourceURL      *string              `json:"source_url"`
		PostedDate     *string              `json:"posted_date"`
		DeadlineDate   *string              `json:"deadline_date"`
		Application    *ApplicationResponse `json:"application,omitempty"`
	} `json:"job"`
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
		Title:          j.Title,
		Location:       j.Location,
		EmploymentType: j.EmploymentType,
	}

	// 企業情報
	if j.Company != nil {
		companyResp := j.Company.ToResponse()
		resp.Company = &companyResp
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
	CompanyID      *int    `form:"company_id"`
	Location       *string `form:"location" validate:"omitempty,max=100"`
	SalaryMin      *int    `form:"salary_min" validate:"omitempty,min=0"`
	EmploymentType *string `form:"employment_type" validate:"omitempty,oneof=正社員 契約社員 業務委託"`
	RemoteOption   *string `form:"remote_option" validate:"omitempty,oneof=リモート可 ハイブリッド 出社必須"`
	AppliedStatus  *string `form:"applied_status" validate:"omitempty,oneof=not_applied applied all"`
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
