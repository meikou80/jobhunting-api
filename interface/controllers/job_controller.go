package controllers

import (
	"net/http"
	"strconv"

	"jobhunting-api/domain/models"
	"jobhunting-api/usecase/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type JobController struct {
	jobService *services.JobService
	validator  *validator.Validate
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{
		jobService: jobService,
		validator:  validator.New(),
	}
}

// CreateJob 求人登録
// @Summary 求人情報を登録
// @Description 新しい求人情報を登録します。重複チェックも同時に行います。
// @Tags jobs
// @Accept json
// @Produce json
// @Param job body models.JobRequest true "求人情報"
// @Success 201 {object} map[string]interface{} "登録成功"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 409 {object} map[string]interface{} "重複データ検出"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /jobs [post]
func (c *JobController) CreateJob(ctx *gin.Context) {
	var req models.JobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	job, duplicates, err := c.jobService.CreateWithDuplicateCheck(&req)
	if err != nil {
		if err.Error() == "duplicate job found" {
			ctx.JSON(http.StatusConflict, gin.H{
				"error":      "Duplicate job detected",
				"duplicates": duplicates,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create job",
			"details": err.Error(),
		})
		return
	}

	response := gin.H{
		"message": "Job created successfully",
		"job":     job.ToResponse(),
	}

	if len(duplicates) > 0 {
		response["duplicates_warning"] = duplicates
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetJobs 求人一覧取得
// @Summary 求人一覧を取得
// @Description フィルタ条件に基づいて求人一覧を取得します
// @Tags jobs
// @Accept json
// @Produce json
// @Param keyword query string false "キーワード検索"
// @Param platform query string false "プラットフォーム"
// @Param company_name query string false "会社名"
// @Param location query string false "勤務地"
// @Param salary_min query int false "最低年収"
// @Param employment_type query string false "雇用形態"
// @Param remote_option query string false "リモート可否"
// @Param status query string false "応募状況"
// @Param show_duplicates query bool false "重複求人も表示"
// @Param page query int false "ページ番号" default(1)
// @Param limit query int false "取得件数" default(20)
// @Success 200 {object} models.JobListResponse "求人一覧"
// @Failure 400 {object} map[string]interface{} "リクエストエラー"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /jobs [get]
func (c *JobController) GetJobs(ctx *gin.Context) {
	var filter models.JobFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}
	if err := c.validator.Struct(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	if filter.Page == 0 {
		filter.Page = 1
	}

	if filter.Limit == 0 {
		filter.Limit = 20
	}

	jobs, total, err := c.jobService.GetJobs(&filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch jobs",
			"details": err.Error(),
		})
		return
	}

	response := models.JobListResponse{
		Jobs:  make([]models.JobResponse, len(jobs)),
		Total: total,
		Page:  filter.Page,
	}

	for i, job := range jobs {
		response.Jobs[i] = job.ToResponse()
	}

	ctx.JSON(http.StatusOK, response)
}

// GetJobDetail 求人詳細取得
// @Summary 求人詳細を取得
// @Description 指定されたIDの求人詳細を取得します（重複求人情報含む）
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "求人ID"
// @Success 200 {object} models.JobDetailResponse "求人詳細"
// @Failure 400 {object} map[string]interface{} "無効なID"
// @Failure 404 {object} map[string]interface{} "求人が見つからない"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /jobs/{id} [get]
func (c *JobController) GetJobDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}
	job, err := c.jobService.GetJobByID(id)
	if err != nil {
		if err.Error() == "job not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch job",
			"details": err.Error(),
		})
		return
	}
	duplicates, err := c.jobService.GetDuplicateJobs(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch duplicate jobs",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, c.buildJobDetailResponse(job, duplicates))
}

// 求人更新API
// PUT api/v1/jobs/:id
func (c *JobController) UpdateJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	var req models.JobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := c.validator.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "validation failed",
			"details": err.Error(),
		})
		return
	}

	job, err := c.jobService.UpdateJob(id, &req)
	if err != nil {
		if err.Error() == "job not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update job",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "job updated successfully",
		"job":     job.ToResponse(),
	})
}

// 求人削除API
// DELETE jobs/{id}
func (c *JobController) DeleteJob(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid job ID",
		})
		return
	}

	err = c.jobService.DeleteJob(id)
	if err != nil {
		if err.Error() == "job not found" {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Job not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete job",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Job deleted successfully",
	})
}

// CheckDuplicates 重複チェック
// @Summary 求人の重複をチェック
// @Description 求人情報から潜在的な重複を検出します
// @Tags jobs
// @Accept json
// @Produce json
// @Param job body models.JobRequest true "チェックする求人情報"
// @Success 200 {object} map[string]interface{} "重複チェック結果"
// @Failure 400 {object} map[string]interface{} "バリデーションエラー"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /jobs/check-duplicates [post]
func (c *JobController) CheckDuplicates(ctx *gin.Context) {
	var req models.JobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	duplicates, err := c.jobService.FindDuplicates(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check duplicates",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"has_duplicates":  len(duplicates) > 0,
		"duplicate_count": len(duplicates),
		"duplicates":      duplicates,
	})
}

// GetPlatformStats プラットフォーム別統計
// @Summary プラットフォーム別統計を取得
// @Description 各プラットフォームの求人数や応募状況の統計を取得します
// @Tags jobs
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "プラットフォーム統計"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /jobs/stats [get]
func (c *JobController) GetPlatformStats(ctx *gin.Context) {
	stats, err := c.jobService.GetPlatformStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch platform stats",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}

// GetDashboard ダッシュボード情報
// @Summary ダッシュボード情報を取得
// @Description 求人統合管理のサマリー情報を取得します
// @Tags jobs
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "ダッシュボード情報"
// @Failure 500 {object} map[string]interface{} "サーバーエラー"
// @Security BearerAuth
// @Router /jobs/dashboard [get]
func (c *JobController) GetDashboard(ctx *gin.Context) {
	dashboard, err := c.jobService.GetDashboard()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch dashboard data",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dashboard)
}

// buildJobDetailResponse 求人詳細レスポンス構築
func (c *JobController) buildJobDetailResponse(job *models.Job, duplicates []models.DuplicateJobInfo) models.JobDetailResponse {
	var response models.JobDetailResponse

	response.Job.ID = job.ID
	response.Job.CompanyName = job.CompanyName
	response.Job.PositionTitle = job.PositionTitle
	response.Job.Description = job.Description
	response.Job.Requirements = job.Requirements
	response.Job.SourcePlatform = job.SourcePlatform
	response.Job.ExternalID = job.ExternalID
	response.Job.SourceURL = job.SourceURL
	response.Job.SalaryMin = job.SalaryMin
	response.Job.SalaryMax = job.SalaryMax
	response.Job.Location = job.Location
	response.Job.EmploymentType = job.EmploymentType
	response.Job.RemoteOption = job.RemoteOption
	response.Job.Status = job.Status
	response.Job.Priority = job.Priority
	response.Job.PersonalNotes = job.PersonalNotes
	response.Job.DuplicateJobs = duplicates

	if job.PostedDate != nil {
		postedDate := job.PostedDate.Format("2006-01-02")
		response.Job.PostedDate = &postedDate
	}

	if job.DeadlineDate != nil {
		deadlineDate := job.DeadlineDate.Format("2006-01-02")
		response.Job.DeadlineDate = &deadlineDate
	}

	return response
}
