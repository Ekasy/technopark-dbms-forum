package repository

import (
	"context"
	"database/sql"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/posts"
	"log"
	"regexp"
)

type PostRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewPostRepository(db *sql.DB) posts.PostRepository {
	return &PostRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (pr *PostRepository) SelectFormSlugByThread(slug string, id int64) (string, int64, error) {
	var queryStr string
	var row *sql.Row
	if len(slug) == 0 {
		queryStr = "SELECT forum, id FROM threads WHERE id = $1;"
		row = pr.db.QueryRow(queryStr, id)
	} else {
		queryStr = "SELECT forum, id FROM threads WHERE slug = $1;"
		row = pr.db.QueryRow(queryStr, slug)
	}

	var forumSlug string
	var threadId int64
	err := row.Scan(&forumSlug, &threadId)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return "", 0, myerr.ThreadNotExists
		}
		pr.logger.Println(err.Error())
		return "", 0, myerr.InternalDbError
	}
	return forumSlug, threadId, nil
}

func (pr *PostRepository) CreatePost(inputPost *models.PostInput, dt string, forumSlug string, threadId int64) (*models.Post, error) {
	tx, err := pr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		pr.logger.Println(err.Error())
	}

	var row *sql.Row
	queryStr := `
		INSERT INTO posts (message, forum, thread, created, author%s) 
		VALUES ($1, $2, $3, $4,
			COALESCE((SELECT nickname FROM users WHERE nickname = $5), $5)
			%s
		)
		RETURNING id, message, forum, thread, created, author, parent, isEdited;`
	if inputPost.Parent == 0 {
		queryStr = fmt.Sprintf(queryStr, "", "")
		row = tx.QueryRow(queryStr, inputPost.Message, forumSlug, threadId, dt, inputPost.Author)
	} else {
		queryStr = fmt.Sprintf(queryStr, ", parent", ", COALESCE((SELECT id FROM posts WHERE id = $6), $6)")
		row = tx.QueryRow(queryStr, inputPost.Message, forumSlug, threadId, dt, inputPost.Author, inputPost.Parent)
	}

	post := &models.Post{}
	err = row.Scan(&post.Id, &post.Message, &post.Forum, &post.Thread, &post.Created, &post.Author, &post.Parent, &post.IsEdited)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}
		fmt.Println(err.Error())
		res, _ := regexp.Match(".*posts_parent_fkey.*", []byte(err.Error()))
		if res {
			return nil, myerr.ParentNotExist
		}

		pr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return nil, myerr.CommitError
	}
	return post, nil
}
