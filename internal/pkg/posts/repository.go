package posts

import "forum/internal/models"

type PostRepository interface {
	SelectFormSlugByThread(slug string, id int64) (string, int64, error)
	CreatePost(inputPost *models.PostInput, dt string, forumSlug string, threadId int64) (*models.Post, error)
	CreatePosts(inputPost []*models.PostInput, dt string, forumSlug string, threadId int64) ([]*models.Post, error)
	CreatePostsAsync(inputPost []*models.PostInput, dt string, forumSlug string, threadId int64) ([]*models.Post, error)
	SelectThreadsBySort(tq *models.ThreadsQuery) ([]*models.Post, error)
	SelectThread(id int64, slug string) (int64, error)
	SelectPost(id int64) (*models.Post, error)
	SelectUser(nickname string) (*models.User, error)
	SelectThreadById(id int64) (*models.Thread, error)
	SelectForum(slug string) (*models.Forum, error)
	UpdatePost(postupdate *models.PostUpdate) (*models.Post, error)
}
