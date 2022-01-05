package posts

import "forum/internal/models"

type PostRepository interface {
	SelectFormSlugByThread(slug string, id int64) (string, int64, error)
	CreatePost(inputPost *models.PostInput, dt string, forumSlug string, threadId int64) (*models.Post, error)
}
