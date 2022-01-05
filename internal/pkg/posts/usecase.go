package posts

import "forum/internal/models"

type PostUsecase interface {
	CreatePostsBySlugOrId(slug string, id int64, postsInput []*models.PostInput) ([]*models.Post, error)
}
