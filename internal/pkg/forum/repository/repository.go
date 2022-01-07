package repository

import (
	"context"
	"database/sql"
	"fmt"
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
		`INSERT INTO forum (slug, title, author) VALUES ($1, $2, COALESCE((SELECT nickname FROM users WHERE nickname = $3), $3)) RETURNING slug, title, author, posts, threads;`,
		forum.Slug, forum.Title, forum.User,
	)

	err = row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts, &forum.Threads)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return myerr.RollbackError
		}

		res, _ := regexp.Match(".*forum_pkey.*", []byte(err.Error()))
		if res {
			return myerr.ForumAlreadyExist
		}
		res, _ = regexp.Match(".*forum_author_fkey.*", []byte(err.Error()))
		if res {
			return myerr.UserNotExist
		}

		fr.logger.Printf(err.Error())
		return myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return myerr.CommitError
	}

	return nil
}

func (fr *ForumRepository) SelectForum(slug string) (*models.Forum, error) {
	row := fr.db.QueryRow(
		`SELECT slug, title, author, posts, threads FROM forum WHERE slug = $1`,
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

func (fr *ForumRepository) SelectUsers(fv *models.ForumUsersQuery) ([]*models.User, error) {
	queryStr := `
					SELECT nickname, fullname, about, email
					FROM forum_users
					WHERE forum LIKE $1 %s
					ORDER BY nickname %s
					LIMIT $2;
				`
	var rows *sql.Rows
	var err error
	if fv.Since != "" {
		queryStr = fmt.Sprintf(queryStr, fmt.Sprintf(`AND nickname %s $3`, fv.Sign), fv.Sorting)
		rows, err = fr.db.Query(queryStr, fv.ForumSlug, fv.Limit, fv.Since)
	} else {
		queryStr = fmt.Sprintf(queryStr, "", fv.Sorting)
		rows, err = fr.db.Query(queryStr, fv.ForumSlug, fv.Limit)
	}
	if err != nil {
		fr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			fr.logger.Println(err.Error())
			return nil, myerr.InternalDbError
		}

		users = append(users, user)
	}

	return users, nil
}
