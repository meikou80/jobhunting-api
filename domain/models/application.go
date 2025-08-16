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
	ID     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID string `json:"user_id" gorm:"type:uuid;not null;index" validate:"required"`
	JobID  int    `json:"job_id" gorm:"not null;index" validate:"required"`

	// 応募経路情報（重要）
	AppliedVia  string    `json:"applied_via" gorm:"not null;size:50" validate:"required,max=50"`
	AppliedDate time.Time `json:"applied_date" gorm:"not null;index"`

	// 応募状況
	Status       ApplicationStatus `json:"status" gorm:"not null;default:applied;size:50" validate:"required,oneof=applied screening interview final offer rejected withdrawn"`
	CurrentStage *string           `json:"current_stage" gorm:"size:100"`

	// 活動記録
	ActivityLogJSON *string      `json:"-" gorm:"type:text"`              // DB格納用JSON文字列
	ActivityLog     *ActivityLog `json:"activity_log,omitempty" gorm:"-"` // Go構造体
	Notes           *string      `json:"notes" gorm:"type:text"`
	NextAction      *string      `json:"next_action" gorm:"type:text"`
	NextActionDate  *time.Time   `json:"next_action_date" gorm:"index"`

	// 結果
	RejectionReason *string `json:"rejection_reason" gorm:"type:text"`
	OfferDetails    *string `json:"offer_details" gorm:"type:text"`

	// リレーション
	Job *Job `json:"job,omitempty" gorm:"foreignKey:JobID"`

	// メタデータ
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// ApplicationRequest 応募登録・更新リクエスト
type ApplicationRequest struct {
	JobID          int                `json:"job_id" validate:"required"`
	AppliedVia     string             `json:"applied_via" validate:"required,max=50"`
	AppliedDate    string             `json:"applied_date" validate:"required"`
	Status         ApplicationStatus  `json:"status" validate:"required,oneof=applied screening interview final offer rejected withdrawn"`
	CurrentStage   *string            `json:"current_stage" validate:"omitempty,max=100"`
	Notes          *string            `json:"notes"`
	NextAction     *string            `json:"next_action"`
	NextActionDate *string            `json:"next_action_date" validate:"omitempty"`
	ActivityLog    *ActivityLogUpdate `json:"activity_log,omitempty"`
}

// ActivityLogUpdate 活動履歴更新用
type ActivityLogUpdate struct {
	Events []ActivityEvent `json:"events"`
}

// ApplicationResponse 応募レスポンス（API用）
type ApplicationResponse struct {
	ID             int          `json:"id"`
	Job            *JobResponse `json:"job,omitempty"`
	Status         string       `json:"status"`
	StatusDisplay  string       `json:"status_display"`
	AppliedVia     string       `json:"applied_via"`
	AppliedDate    string       `json:"applied_date"`
	CurrentStage   *string      `json:"current_stage"`
	Notes          *string      `json:"notes"`
	NextAction     *string      `json:"next_action"`
	NextActionDate *string      `json:"next_action_date"`
	UpdatedAt      string       `json:"updated_at"`
}

// ApplicationDetailResponse 応募詳細レスポンス
type ApplicationDetailResponse struct {
	Application struct {
		ID             int          `json:"id"`
		Job            *JobResponse `json:"job"`
		Status         string       `json:"status"`
		AppliedVia     string       `json:"applied_via"`
		AppliedDate    string       `json:"applied_date"`
		CurrentStage   *string      `json:"current_stage"`
		Notes          *string      `json:"notes"`
		NextAction     *string      `json:"next_action"`
		NextActionDate *string      `json:"next_action_date"`
		ActivityLog    *ActivityLog `json:"activity_log"`
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
		ID:            a.ID,
		Status:        string(a.Status),
		StatusDisplay: a.GetStatusDisplay(),
		AppliedVia:    a.AppliedVia,
		AppliedDate:   a.AppliedDate.Format("2006-01-02"),
		CurrentStage:  a.CurrentStage,
		Notes:         a.Notes,
		NextAction:    a.NextAction,
		UpdatedAt:     a.UpdatedAt.Format(time.RFC3339),
	}

	// 次のアクション日
	if a.NextActionDate != nil {
		nextActionDate := a.NextActionDate.Format("2006-01-02")
		resp.NextActionDate = &nextActionDate
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
	Status      *string `form:"status" validate:"omitempty,oneof=applied screening interview final offer rejected withdrawn"`
	AppliedVia  *string `form:"applied_via" validate:"omitempty,max=50"`
	CompanyName *string `form:"company_name" validate:"omitempty,max=100"`
	Sort        *string `form:"sort" validate:"omitempty,oneof=applied_date_desc status_asc updated_desc"`
	Page        int     `form:"page" validate:"min=1" default:"1"`
	Limit       int     `form:"limit" validate:"min=1,max=100" default:"20"`
}

// Value SQLドライバー用
func (a Application) Value() (driver.Value, error) {
	return a.ID, nil
}
