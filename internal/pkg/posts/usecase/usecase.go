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
	forumSlug, threadId, err := pu.repo.SelectFormSlugByThread(slug, id)
	switch err {
	case nil:
		// skip this state
	case myerr.ThreadNotExists:
		return nil, myerr.ThreadNotExists
	default:
		return nil, err
	}

	if len(postsInput) == 0 {
		return make([]*models.Post, 0), nil
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

func (pu *PostUsecase) GetInfo(pq *models.PostQuery) (map[string]interface{}, error) {
	post, err := pu.repo.SelectPost(pq.PostId)
	if err != nil {
		return nil, err
	}

	info := make(map[string]interface{})
	info["post"] = post

	for _, val := range pq.Related {
		switch val {
		case "user":
			user, err := pu.repo.SelectUser(post.Author)
			if err != nil {
				return nil, err
			}
			info["author"] = user
		case "forum":
			forum, err := pu.repo.SelectForum(post.Forum)
			if err != nil {
				return nil, err
			}
			info["forum"] = forum
		case "thread":
			thread, err := pu.repo.SelectThreadById(post.Thread)
			if err != nil {
				return nil, err
			}
			info["thread"] = thread
		}
	}
	return info, nil
}

func (pu *PostUsecase) UpdatePost(postupdate *models.PostUpdate) (*models.Post, error) {
	post, err := pu.repo.UpdatePost(postupdate)
	return post, err
}
