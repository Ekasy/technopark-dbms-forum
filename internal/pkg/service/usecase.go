package service

import "forum/internal/models"

type ServiceUsecase interface {
	GetServiceStatus() (*models.Service, error)
	ClearService() error
}
