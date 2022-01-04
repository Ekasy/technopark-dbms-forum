package repository

import (
	"context"
	"database/sql"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/threads"
	"log"
	"regexp"
)

type ThreadRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewThreadRepository(db *sql.DB) threads.ThreadRepository {
	return &ThreadRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (tr *ThreadRepository) InsertThread(thread *models.Thread) error {
	tx, err := tr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return myerr.InternalDbError
	}

	row := tx.QueryRow(
		`INSERT INTO threads (title, message, slug, author, forum, created) 
		 VALUES ($1, $2, $3,
			COALESCE((SELECT nickname FROM users WHERE nickname = $4), $4),
			COALESCE((SELECT slug FROM forum WHERE slug = $5), $5),
			$6
		 )
		 RETURNING id, title, author, forum, message, votes, slug, created;`,
		thread.Title, thread.Message, thread.Slug, thread.Author, thread.Forum, thread.Created,
	)

	err = row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return myerr.RollbackError
		}

		res, _ := regexp.Match(".*threads_pkey.*", []byte(err.Error()))
		if res {
			return myerr.ThreadAlreadyExist
		}

		res, _ = regexp.Match(".*threads_author_fkey.*", []byte(err.Error()))
		if res {
			return myerr.AuthorNotExist
		}

		res, _ = regexp.Match(".*threads_forum_fkey.*", []byte(err.Error()))
		if res {
			return myerr.ForumNotExist
		}

		tr.logger.Printf(err.Error())
		return myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return myerr.CommitError
	}
	return nil
}

func (tr *ThreadRepository) SelectThreadBySlug(slug string) (*models.Thread, error) {
	thread := &models.Thread{}
	row := tr.db.QueryRow(
		"SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1",
		slug,
	)
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.NoRows
		}
		return nil, myerr.InternalDbError
	}
	return thread, nil
}
