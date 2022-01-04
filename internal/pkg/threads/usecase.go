package threads

import "forum/internal/models"

type ThreadUsecase interface {
	CreateThread(thread *models.Thread) (*models.Thread, error)
}
