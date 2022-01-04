package threads

import "forum/internal/models"

type ThreadRepository interface {
	InsertThread(thread *models.Thread) error
	SelectThreadBySlug(slug string) (*models.Thread, error)
}
