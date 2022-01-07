package usecase

import (
	"forum/internal/models"
	"forum/internal/pkg/service"
)

type ServiceUsecase struct {
	repo service.ServiceRepository
}

func NewServiceUsecase(repo service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		repo: repo,
	}
}

func (su *ServiceUsecase) GetServiceStatus() (*models.Service, error) {
	srvc, err := su.repo.SelectServiceStatus()
	return srvc, err
}

func (su *ServiceUsecase) ClearService() error {
	err := su.repo.ClearService()
	return err
}
