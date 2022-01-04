package forum

import "forum/internal/models"

type ForumUsecase interface {
	CreateForum(forum *models.Forum) (*models.Forum, error)
	GetForum(slug string) (*models.Forum, error)
}
