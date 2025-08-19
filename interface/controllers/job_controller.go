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

// 求人登録API
// POST /api/v1/jobs
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

// 求人一覧取得API
// GET /api/v1/jobs
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

// 求人詳細取得API
// GET /api/v1/jobs/:id
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

// 求人重複チェックAPI
// GET /api/v1/jobs/duplicates
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

// プラットフォーム統計API
// GET /api/v1/jobs/stats
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

// ダッシュボードAPI
// GET /api/v1/jobs/dashboard
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
