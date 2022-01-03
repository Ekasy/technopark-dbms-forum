package usecase

import (
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/user"
	"log"
)

type UserUsecase struct {
	repo   user.UserRepository
	logger *log.Logger
}

func NewUserUsecase(repo user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		repo:   repo,
		logger: log.Default(),
	}
}

func (uu *UserUsecase) GetUser(nickname string) (*models.User, error) {
	user, err := uu.repo.SelectUser(nickname)
	return user, err
}

func (uu *UserUsecase) CreateUser(user *models.User) ([]*models.User, bool, error) {
	err := uu.repo.InsertUser(user)
	users := make([]*models.User, 0)

	switch err {
	case nil:
		users = append(users, user)
		return append(users, user), true, err
	case myerr.NicknameAlreadyExist:
		users, err = uu.repo.SelectUsersIfExists(user.Nickname, user.Email)
		return users, false, err
	case myerr.EmailAlreadyExist:
		users, err = uu.repo.SelectUsersIfExists(user.Nickname, user.Email)
		return users, false, err
	default:
		return nil, false, err
	}
}

func (uu *UserUsecase) UpdateUser(user *models.User) (*models.User, error) {
	err := uu.repo.UpdateUser(user)
	return user, err
}
