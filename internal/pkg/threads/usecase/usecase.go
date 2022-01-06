package usecase

import (
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/threads"
)

type ThreadUsecase struct {
	repo threads.ThreadRepository
}

func NewThreadUsecase(repo threads.ThreadRepository) threads.ThreadUsecase {
	return &ThreadUsecase{
		repo: repo,
	}
}

func (tu *ThreadUsecase) CreateThread(thread *models.Thread) (*models.Thread, error) {
	err := tu.repo.InsertThread(thread)
	switch err {
	case nil:
		return thread, nil
	case myerr.AuthorNotExist:
		return nil, myerr.AuthorNotExist
	case myerr.ForumNotExist:
		return nil, myerr.ForumNotExist
	case myerr.ThreadAlreadyExist:
		thread, err = tu.repo.SelectThreadBySlug(thread.Slug)
		if err != nil {
			return nil, err
		}
		return thread, myerr.ThreadAlreadyExist
	default:
		return nil, err
	}
}

func (tu *ThreadUsecase) GetThreadsByForum(tv *models.ThreadsVars) ([]*models.Thread, error) {
	threads, err := tu.repo.SelectThreadsByForum(tv)
	return threads, err
}

func (tu *ThreadUsecase) GetUsersByForum(tv *models.ThreadsVars) ([]*models.Thread, error) {
	threads, err := tu.repo.SelectUsersByForum(tv)
	return threads, err
}

func (tu *ThreadUsecase) GetThread(slug string, id int64) (*models.Thread, error) {
	thread, err := tu.repo.SelectThread(slug, id)
	return thread, err
}
