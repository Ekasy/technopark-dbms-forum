package service

import "forum/internal/models"

type ServiceRepository interface {
	SelectServiceStatus() (*models.Service, error)
	ClearService() error
}
