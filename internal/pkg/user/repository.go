package user

import "forum/internal/models"

type UserRepository interface {
	InsertUser(user *models.User) error
	UpdateUser(user *models.User) error
	SelectUser(nickname string) (*models.User, error)
	SelectUsersIfExists(nickname string, email string) ([]*models.User, error)
}
