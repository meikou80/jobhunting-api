package services

import (
	"context"
	"fmt"

	"jobhunting-api/domain/models"
	"jobhunting-api/domain/repositories"
)

type CompanyService interface {
	CreateCompany(ctx context.Context, req *models.CompanyRequest) (*models.Company, error)
	GetCompany(ctx context.Context, id int) (*models.Company, error)
	UpdateCompany(ctx context.Context, id int, req *models.CompanyRequest) (*models.Company, error)
	DeleteCompany(ctx context.Context, id int) error
	ListCompanies(ctx context.Context, filter *models.CompanyFilter) (*models.CompanyListResponse, error)
	SearchCompanies(ctx context.Context, query string) ([]*models.Company, error)
}

type companyService struct {
	companyRepo repositories.CompanyRepository
}

func NewCompanyService(companyRepo repositories.CompanyRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
	}
}

func (s *companyService) CreateCompany(ctx context.Context, req *models.CompanyRequest) (*models.Company, error) {
	company := &models.Company{
		Name:         req.Name,
		Industry:     req.Industry,
		SizeCategory: req.SizeCategory,
		Location:     req.Location,
		WebsiteURL:   req.WebsiteURL,
		Notes:        req.Notes,
		Rating:       req.Rating,
	}

	existing, err := s.companyRepo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("company with name '%s' already exists", req.Name)
	}

	if err := s.companyRepo.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	return company, nil

}

func (s *companyService) GetCompany(ctx context.Context, id int) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if company == nil {
		return nil, fmt.Errorf("company with id %d not found", id)
	}

	return company, nil
}

func (s *companyService) UpdateCompany(ctx context.Context, id int, req *models.CompanyRequest) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if company == nil {
		return nil, fmt.Errorf("company with id %d not found", id)
	}

	if req.Name != company.Name {
		existing, err := s.companyRepo.GetByName(ctx, req.Name)
		if err == nil && existing != nil && existing.ID != id {
			return nil, fmt.Errorf("company with name '%s' already exists", req.Name)
		}
	}

	company.Name = req.Name
	company.Industry = req.Industry
	company.SizeCategory = req.SizeCategory
	company.Location = req.Location
	company.WebsiteURL = req.WebsiteURL
	company.Notes = req.Notes
	company.Rating = req.Rating

	if err := s.companyRepo.Update(ctx, company); err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	return company, nil
}

func (s *companyService) DeleteCompany(ctx context.Context, id int) error {
	company, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	if company == nil {
		return fmt.Errorf("company with id %d not found", id)
	}

	if err := s.companyRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	return nil
}

func (s *companyService) ListCompanies(ctx context.Context, filter *models.CompanyFilter) (*models.CompanyListResponse, error) {
	companies, total, err := s.companyRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}

	companyResponses := make([]models.CompanyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = company.ToResponse()
	}

	return &models.CompanyListResponse{
		Companies: companyResponses,
		Total:     total,
		Page:      filter.Page,
	}, nil
}

func (s *companyService) SearchCompanies(ctx context.Context, keyword string) ([]*models.Company, error) {
	companies, err := s.companyRepo.Search(ctx, keyword)
	if err != nil {
		return nil, fmt.Errorf("failed to search companies: %w", err)
	}

	return companies, nil
}
