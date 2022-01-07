package forum

import "forum/internal/models"

type ForumRepository interface {
	InsertForum(forum *models.Forum) error
	SelectForum(slug string) (*models.Forum, error)
	SelectUsers(fv *models.ForumUsersQuery) ([]*models.User, error)
}
