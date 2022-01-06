package usecase

import (
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/posts"
	"time"
)

type PostUsecase struct {
	repo posts.PostRepository
}

func NewPostUsecase(repo posts.PostRepository) posts.PostUsecase {
	return &PostUsecase{
		repo: repo,
	}
}

func (pu *PostUsecase) CreatePostsBySlugOrId(slug string, id int64, postsInput []*models.PostInput) ([]*models.Post, error) {
	if len(postsInput) == 0 {
		return make([]*models.Post, 0), nil
	}
	forumSlug, threadId, err := pu.repo.SelectFormSlugByThread(slug, id)
	switch err {
	case nil:
		// skip this state
	case myerr.ThreadNotExists:
		return nil, err
	default:
		return nil, err
	}

	dt := time.Now().Format(models.Layout)
	posts := make([]*models.Post, 0)
	for _, pi := range postsInput {
		post, err := pu.repo.CreatePost(pi, dt, forumSlug, threadId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)

	}

	return posts, nil
}

func (pu *PostUsecase) GetPostsRec(tq *models.ThreadsQuery) ([]*models.Post, error) {
	id, err := pu.repo.SelectThread(tq.ThreadId, tq.ThreadSlug)
	if err != nil {
		return nil, err
	}

	tq.ThreadId = id
	posts, err := pu.repo.SelectThreadsBySort(tq)
	return posts, err
}
