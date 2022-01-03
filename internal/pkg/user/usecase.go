package user

import "forum/internal/models"

type UserUsecase interface {
	GetUser(nickname string) (*models.User, error)
	CreateUser(user *models.User) ([]*models.User, bool, error)
	UpdateUser(user *models.User) (*models.User, error)
}
