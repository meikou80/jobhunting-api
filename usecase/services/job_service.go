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
