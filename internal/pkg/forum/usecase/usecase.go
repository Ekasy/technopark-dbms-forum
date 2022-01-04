package usecase

import (
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/forum"
)

type ForumUsecase struct {
	repo forum.ForumRepository
}

func NewForumUsecase(repo forum.ForumRepository) forum.ForumUsecase {
	return &ForumUsecase{
		repo: repo,
	}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) (*models.Forum, error) {
	err := fu.repo.InsertForum(forum)
	if err == myerr.ForumAlreadyExist {
		forum, err = fu.repo.SelectForum(forum.Slug)
		if err == nil {
			err = myerr.ForumAlreadyExist
		}
	}
	return forum, err
}

func (fu *ForumUsecase) GetForum(slug string) (*models.Forum, error) {
	forum, err := fu.repo.SelectForum(slug)
	return forum, err
}
