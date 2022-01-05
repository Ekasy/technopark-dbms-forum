package threads

import "forum/internal/models"

type ThreadUsecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
	GetThreadsByForum(tv *models.ThreadsVars) ([]*models.Thread, error)
	GetUsersByForum(tv *models.ThreadsVars) ([]*models.Thread, error)
}
