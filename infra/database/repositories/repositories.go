package repositories

import (
	"gorm.io/gorm"

	"jobhunting-api/domain/repositories"
)

// Repositories リポジトリの集約
type Repositories struct {
	Company     repositories.CompanyRepository
	Job         repositories.JobRepository
	Application repositories.ApplicationRepository
	User        repositories.UserRepository
}

// NewRepositories 全リポジトリのコンストラクタ
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Company:     NewCompanyRepository(db),
		Job:         NewJobRepository(db),
		Application: NewApplicationRepository(db),
		User:        NewUserRepository(db),
	}
}
