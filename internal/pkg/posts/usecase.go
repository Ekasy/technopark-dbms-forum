package posts

import "forum/internal/models"

type PostUsecase interface {
	CreatePostsBySlugOrId(slug string, id int64, postsInput []*models.PostInput) ([]*models.Post, error)
	GetPostsRec(tq *models.ThreadsQuery) ([]*models.Post, error)
	GetInfo(pq *models.PostQuery) (map[string]interface{}, error)
	UpdatePost(pu *models.PostUpdate) (*models.Post, error)
}
