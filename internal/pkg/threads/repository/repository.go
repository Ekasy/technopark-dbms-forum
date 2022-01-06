package repository

import (
	"context"
	"database/sql"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/threads"
	"log"
	"regexp"
	"time"
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

		res, _ = regexp.Match(".*pindex_threads_slug.*", []byte(err.Error()))
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

func (tr *ThreadRepository) SelectThreadsByForum(tv *models.ThreadsVars) ([]*models.Thread, error) {
	row := tr.db.QueryRow(
		"SELECT slug FROM forum WHERE slug = $1",
		tv.ForumSlug,
	)
	buf := ""
	err := row.Scan(&buf)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.NoRows
		}
		tr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	queryStr := `
					SELECT id, title, author, forum, message, votes, slug, created
					FROM threads
					WHERE forum LIKE $1 %s
					ORDER BY created %s
					LIMIT $2;
				`
	var rows *sql.Rows
	if tv.Since != "" {
		queryStr = fmt.Sprintf(queryStr, fmt.Sprintf(`AND created::timestamp %s $3::timestamp with time zone`, tv.Sign), tv.Sorting)
		rows, err = tr.db.Query(queryStr, tv.ForumSlug, tv.Limit, tv.Since)
	} else {
		queryStr = fmt.Sprintf(queryStr, "", tv.Sorting)
		rows, err = tr.db.Query(queryStr, tv.ForumSlug, tv.Limit)
	}
	if err != nil {
		tr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		thread := &models.Thread{}
		t := &time.Time{}
		err = rows.Scan(
			&thread.Id, &thread.Title, &thread.Author, &thread.Forum,
			&thread.Message, &thread.Votes, &thread.Slug, &t)
		if err != nil {
			tr.logger.Println(err.Error())
			return nil, myerr.InternalDbError
		}

		thread.Created = t.Format(models.Layout)
		threads = append(threads, thread)
	}

	return threads, nil
}

func (tr *ThreadRepository) SelectUsersByForum(tv *models.ThreadsVars) ([]*models.Thread, error) {
	row := tr.db.QueryRow(
		"SELECT slug FROM forum WHERE slug = $1",
		tv.ForumSlug,
	)
	buf := ""
	err := row.Scan(&buf)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.NoRows
		}
		tr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	return nil, nil
}

func (tr *ThreadRepository) SelectThread(slug string, id int64) (*models.Thread, error) {
	thread := &models.Thread{}
	row := tr.db.QueryRow(
		"SELECT id, title, author, forum, message, votes, slug, created from threads WHERE 0 = $1 AND slug LIKE $2 OR $2 LIKE '' AND id = $1",
		id, slug,
	)
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.ThreadNotExists
		}
		tr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}
	return thread, nil
}

func (tr *ThreadRepository) UpdateThread(threadUpdate *models.ThreadUpdate) (*models.Thread, error) {
	tx, err := tr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, myerr.InternalDbError
	}

	thread := &models.Thread{}
	row := tx.QueryRow(
		`UPDATE threads SET 
			title = $1, 
			message = $2 
		 WHERE id = (SELECT id from threads WHERE 0 = $3 AND slug LIKE $4 OR $4 LIKE '' AND id = $3)
		 RETURNING id, title, author, forum, message, votes, slug, created;`,
		threadUpdate.Title, threadUpdate.Message, threadUpdate.Id, threadUpdate.Slug)

	err = row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}

		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.ThreadNotExists
		}

		tr.logger.Printf(err.Error())
		return nil, myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return nil, myerr.CommitError
	}
	return thread, nil
}
