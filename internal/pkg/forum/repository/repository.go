package repository

import (
	"context"
	"database/sql"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/forum"
	"log"
	"regexp"
)

type ForumRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewForumRepository(db *sql.DB) forum.ForumRepository {
	return &ForumRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (fr *ForumRepository) InsertForum(forum *models.Forum) error {
	tx, err := fr.db.BeginTx(context.Background(), nil)
	if err != nil {
		return myerr.InternalDbError
	}

	row := tx.QueryRowContext(
		context.Background(),
		`INSERT INTO forum (slug, title, "user") VALUES ($1, $2, $3) RETURNING slug, title, user, posts, threads;`,
		forum.Slug, forum.Title, forum.User,
	)

	err = row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)
	if err != nil {
		res, _ := regexp.Match(".*forum_pkey.*", []byte(err.Error()))
		if res {
			return myerr.ForumAlreadyExist
		}
		res, _ = regexp.Match(".*forum_user_fkey.*", []byte(err.Error()))
		if res {
			return myerr.UserNotExist
		}
	}

	if err != nil {
		fr.logger.Printf(err.Error())
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return myerr.RollbackError
		}
		return myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return myerr.CommitError
	}
	return nil
}

func (fr *ForumRepository) SelectForum(slug string) (*models.Forum, error) {
	row := fr.db.QueryRowContext(
		context.Background(),
		`SELECT slug, title, "user", posts, threads FROM forum WHERE slug = $1`,
		slug,
	)

	forum := &models.Forum{}
	err := row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.NoRows
		}
	}

	if err != nil {
		fr.logger.Printf(err.Error())
		return nil, myerr.InternalDbError
	}

	return forum, nil
}
