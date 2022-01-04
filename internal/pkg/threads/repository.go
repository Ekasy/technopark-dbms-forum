package threads

import "forum/internal/models"

type ThreadRepository interface {
	InsertThread(thread *models.Thread)
}
