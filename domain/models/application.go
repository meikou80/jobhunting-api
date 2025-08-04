package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ApplicationStatus 応募ステータス
type ApplicationStatus string

const (
	StatusInterested ApplicationStatus = "interested" // 興味あり
	StatusApplied    ApplicationStatus = "applied"    // 応募済み
	StatusScreening  ApplicationStatus = "screening"  // 書類選考中
	StatusInterview  ApplicationStatus = "interview"  // 面接中
	StatusFinal      ApplicationStatus = "final"      // 最終面接
	StatusOffer      ApplicationStatus = "offer"      // 内定
	StatusRejected   ApplicationStatus = "rejected"   // 不採用
	StatusWithdrawn  ApplicationStatus = "withdrawn"  // 辞退
)

// ActivityEvent 活動履歴イベント
type ActivityEvent struct {
	Date        string  `json:"date"`
	Type        string  `json:"type"`
	Note        string  `json:"note"`
	Interviewer *string `json:"interviewer,omitempty"`
	Result      *string `json:"result,omitempty"`
}

// ActivityLog 活動履歴
type ActivityLog struct {
	Events []ActivityEvent `json:"events"`
}

// Application 応募記録
type Application struct {
	// 基本情報
	ID        int    `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	JobID     int    `json:"job_id" db:"job_id"`
	CompanyID *int   `json:"company_id" db:"company_id"`

	// 応募情報
	Status      ApplicationStatus `json:"status" db:"status" validate:"required,oneof=interested applied screening interview final offer rejected withdrawn"`
	AppliedDate *time.Time        `json:"applied_date" db:"applied_date"`
	Priority    int               `json:"priority" db:"priority" validate:"min=1,max=5"`

	// 活動記録
	ActivityLogJSON *string      `json:"-" db:"activity_log"`           // DB格納用JSON文字列
	ActivityLog     *ActivityLog `json:"activity_log,omitempty" db:"-"` // Go構造体
	Notes           *string      `json:"notes" db:"notes"`
	InterviewNotes  *string      `json:"interview_notes" db:"interview_notes"`
	CompanyResearch *string      `json:"company_research" db:"company_research"`

	// 結果
	RejectionReason *string `json:"rejection_reason" db:"rejection_reason"`
	OfferDetails    *string `json:"offer_details" db:"offer_details"`

	// リレーション
	Job     *Job     `json:"job,omitempty"`
	Company *Company `json:"company,omitempty"`

	// メタデータ
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ApplicationRequest 応募登録・更新リクエスト
type ApplicationRequest struct {
	JobID       int                `json:"job_id" validate:"required"`
	Status      ApplicationStatus  `json:"status" validate:"required,oneof=interested applied screening interview final offer rejected withdrawn"`
	AppliedDate *string            `json:"applied_date" validate:"omitempty"`
	Priority    int                `json:"priority" validate:"min=1,max=5"`
	Notes       *string            `json:"notes"`
	ActivityLog *ActivityLogUpdate `json:"activity_log,omitempty"`
}

// ActivityLogUpdate 活動履歴更新用
type ActivityLogUpdate struct {
	Events []ActivityEvent `json:"events"`
}

// ApplicationResponse 応募レスポンス（API用）
type ApplicationResponse struct {
	ID              int          `json:"id"`
	Job             *JobResponse `json:"job,omitempty"`
	Status          string       `json:"status"`
	StatusDisplay   string       `json:"status_display"`
	AppliedDate     *string      `json:"applied_date"`
	Priority        int          `json:"priority"`
	PriorityDisplay string       `json:"priority_display"`
	Notes           *string      `json:"notes"`
	NextAction      *string      `json:"next_action"`
	UpdatedAt       string       `json:"updated_at"`
}

// ApplicationDetailResponse 応募詳細レスポンス
type ApplicationDetailResponse struct {
	Application struct {
		ID              int          `json:"id"`
		Job             *JobResponse `json:"job"`
		Status          string       `json:"status"`
		AppliedDate     *string      `json:"applied_date"`
		Priority        int          `json:"priority"`
		Notes           *string      `json:"notes"`
		InterviewNotes  *string      `json:"interview_notes"`
		CompanyResearch *string      `json:"company_research"`
		ActivityLog     *ActivityLog `json:"activity_log"`
	} `json:"application"`
}

// ApplicationListResponse 応募一覧レスポンス
type ApplicationListResponse struct {
	Applications []ApplicationResponse `json:"applications"`
	Summary      struct {
		Total    int            `json:"total"`
		ByStatus map[string]int `json:"by_status"`
	} `json:"summary"`
}

// ToResponse Application → ApplicationResponse 変換
func (a *Application) ToResponse() ApplicationResponse {
	resp := ApplicationResponse{
		ID:              a.ID,
		Status:          string(a.Status),
		StatusDisplay:   a.GetStatusDisplay(),
		Priority:        a.Priority,
		PriorityDisplay: a.GetPriorityDisplay(),
		Notes:           a.Notes,
		UpdatedAt:       a.UpdatedAt.Format(time.RFC3339),
	}

	// 応募日
	if a.AppliedDate != nil {
		appliedDate := a.AppliedDate.Format("2006-01-02")
		resp.AppliedDate = &appliedDate
	}

	// 求人情報
	if a.Job != nil {
		jobResp := a.Job.ToResponse()
		resp.Job = &jobResp
	}

	return resp
}

// GetStatusDisplay ステータス表示名を取得
func (a *Application) GetStatusDisplay() string {
	statusMap := map[ApplicationStatus]string{
		StatusInterested: "興味あり",
		StatusApplied:    "応募済み",
		StatusScreening:  "書類選考中",
		StatusInterview:  "面接中",
		StatusFinal:      "最終面接",
		StatusOffer:      "内定",
		StatusRejected:   "不採用",
		StatusWithdrawn:  "辞退",
	}
	return statusMap[a.Status]
}

// GetPriorityDisplay 優先度表示名を取得
func (a *Application) GetPriorityDisplay() string {
	priorityMap := map[int]string{
		1: "最高",
		2: "高",
		3: "中",
		4: "低",
		5: "最低",
	}
	return priorityMap[a.Priority]
}

// SetActivityLog ActivityLogをJSON文字列に変換してDBに格納
func (a *Application) SetActivityLog(log *ActivityLog) error {
	if log == nil {
		a.ActivityLogJSON = nil
		a.ActivityLog = nil
		return nil
	}

	jsonBytes, err := json.Marshal(log)
	if err != nil {
		return err
	}

	jsonStr := string(jsonBytes)
	a.ActivityLogJSON = &jsonStr
	a.ActivityLog = log
	return nil
}

// GetActivityLog JSON文字列からActivityLogを復元
func (a *Application) GetActivityLog() (*ActivityLog, error) {
	if a.ActivityLogJSON == nil || *a.ActivityLogJSON == "" {
		return nil, nil
	}

	var log ActivityLog
	if err := json.Unmarshal([]byte(*a.ActivityLogJSON), &log); err != nil {
		return nil, err
	}

	return &log, nil
}

// ApplicationFilter 応募検索フィルタ
type ApplicationFilter struct {
	Status    *string `form:"status" validate:"omitempty,oneof=interested applied screening interview final offer rejected withdrawn"`
	Priority  *int    `form:"priority" validate:"omitempty,min=1,max=5"`
	CompanyID *int    `form:"company_id"`
	Sort      *string `form:"sort" validate:"omitempty,oneof=applied_date_desc priority_asc updated_desc"`
	Page      int     `form:"page" validate:"min=1" default:"1"`
	Limit     int     `form:"limit" validate:"min=1,max=100" default:"20"`
}

// Value SQLドライバー用
func (a Application) Value() (driver.Value, error) {
	return a.ID, nil
}
