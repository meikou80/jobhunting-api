package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

type JobService struct {
	jobRepo repositories.JobRepository
}

func NewJobService(jobRepo repositories.JobRepository) *JobService {
	return &JobService{
		jobRepo: jobRepo,
	}
}

func (s *JobService) CreateWithDuplicateCheck(req *models.JobRequest) (*models.Job, []models.DuplicateJobInfo, error) {
	ctx := context.Background()

	job, err := s.convertRequestToJob(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert request: %w", err)
	}

	duplicates, err := s.FindDuplicates(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find duplicates: %w", err)
	}

	if len(duplicates) > 0 {
		return nil, duplicates, errors.New("duplicate job found")
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, nil, fmt.Errorf("failed to create job: %w", err)
	}

	return job, duplicates, nil
}

// GetJobs 求人一覧取得
func (s *JobService) GetJobs(filter *models.JobFilter) ([]*models.Job, int, error) {
	ctx := context.Background()
	return s.jobRepo.List(ctx, filter)
}

// GetJobByID 求人詳細取得
func (s *JobService) GetJobByID(id int) (*models.Job, error) {
	ctx := context.Background()
	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errors.New("job not found")
	}
	return job, nil
}

// UpdateJob 求人更新
func (s *JobService) UpdateJob(id int, req *models.JobRequest) (*models.Job, error) {
	ctx := context.Background()

	existingJob, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingJob == nil {
		return nil, errors.New("job not found")
	}

	updatedJob, err := s.convertRequestToJob(req)
	if err != nil {
		return nil, fmt.Errorf("failed to convert request: %w", err)
	}

	updatedJob.ID = existingJob.ID
	updatedJob.CreatedAt = existingJob.CreatedAt
	updatedJob.UpdatedAt = time.Now()

	if err := s.jobRepo.Update(ctx, updatedJob); err != nil {
		return nil, fmt.Errorf("failed to update job: %w", err)
	}

	return updatedJob, nil
}

// DeleteJob 求人削除
func (s *JobService) DeleteJob(id int) error {
	ctx := context.Background()

	job, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if job == nil {
		return errors.New("job not found")
	}

	if err := s.jobRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	return nil
}

// FindDuplicates 重複求人検索
func (s *JobService) FindDuplicates(req *models.JobRequest) ([]models.DupilcateJobInfo, error) {
	ctx := context.Background()
	var duplicates []models.DuplicateJobInfo

	// 1. 会社名での検索
	companyJobs, err := s.jobRepo.GetByCompanyName(ctx, req.CompanyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get company jobs: %w", err)
	}

	for _, job := range companyJobs {
		confidence := s.calculateDuplicateConfidence(req, job)
		if confidence > 0.7 {
			duplicates = append(duplicates, models.DuplicateJobInfo{
				ID:             job.ID,
				SourcePlatform: job.SourcePlatform,
				ExternalID:     job.ExternalID,
				SourceURL:      job.SourceURL,
				Confidence:     confidence,
			})
		}
	}

	// 2. 外部IDでの完全一致チェック
	if req.ExternalID != nil {
		existingJob, err := s.jobRepo.GetByExternalID(ctx, req.SourcePlatform, *req.ExternalID)
		if err == nil && existingJob != nil {
			found := false
			for _, dup := range duplicates {
				if dup.ID == existingJob.ID {
					found = true
					break
				}
			}
			if !found {
				duplicates = append(duplicates, models.DuplicateJobInfo{
					ID:             existingJob.ID,
					SourcePlatform: existingJob.SourcePlatform,
					ExternalID:     existingJob.ExternalID,
					SourceURL:      existingJob.SourceURL,
					Confidence:     1.0,
				})
			}
		}
	}
	return duplicates, nil
}

// GetDuplicateJobs 指定求人の重複一覧取得
func (s *JobService) GetDuplicateJobs(jobID int) ([]models.DuplicateJobInfo, error) {
	ctx := context.Background()

	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errors.New("job not found")
	}
	req := &models.JobRequest{
		CompanyName:    job.CompanyName,
		PositionTitle:  job.PositionTitle,
		SourcePlatform: job.SourcePlatform,
		ExternalID:     job.ExternalID,
	}

	duplicates, err := s.FindDuplicates(req)
	if err != nil {
		return nil, err
	}

	var filteredDuplicates []models.DuplicateJobInfo
	for _, dup := range duplicates {
		if dup.ID != jobID {
			filteredDuplicates = append(filteredDuplicates, dup)
		}
	}

	return filteredDuplicates, nil

}

// GetPlatformStats プラットフォーム別統計
func (s *JobService) GetPlatformStats() (map[string]interface{}, error) {
	ctx := context.Background()

	platforms := []string("doda", "recruit", "geekly", "levtech", "その他")
	platformStats := make(map[string]int64)

	for _, platform := range platforms {
		count, err = s.jobRepo.CountByPlatform(ctx, platform)
		if err != nil {
			return nil, fmt.Errorf("failed to count jobs for platform %s: %w", platform, err)
		}
		platformStats[platform] = count
	}

	totalJobs, err := s.jobRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total jobs: %w", err)
	}

	stats := map[string]interface{}{
		"total_jobs":     totalJobs,
		"platform_stats": platformStats,
		"created_at":     time.Now().Format("2006-01-02"),
	}

	return stats, nil
}

// GetDashboard ダッシュボード情報取得
func (s *JobService) GetDashboard() (map[string]interface{}, error) {
	ctx := context.Background()

	totalJobs, err := s.jobRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total jobs: %w", err)
	}
	statusStats := make(map[string]int)
	statuses := []string("未対応", "対応中", "対応済み")

	for _, status := range statuses {
		filter := &models.Jobfilter{
			Status: &status,
			Page:   1,
			Limit:  1000,
		}
		jobs, total, err := s.jobRepo.List(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get jobs by status %s: %w", status, err)
		}
		_ = jobs
		statusStats[status] = total
	}

	platformStats, err := s.GetPlatformStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get platform stats: %W", err)
	}

	dashboard := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_jobs":   totalJobs,
			"status_stats": statusStats,
		},
		"platform_stats": platformStats,
		"last_updated":   time.Now().Format("2006-01-02 15:04:05"),
	}

	return dashboard, nil
}
