package repository

import (
	"database/sql"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/service"
	"log"
)

type ServiceRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewServiceRepository(db *sql.DB) service.ServiceRepository {
	return &ServiceRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (sr *ServiceRepository) SelectServiceStatus() (*models.Service, error) {
	srvc := &models.Service{}
	row := sr.db.QueryRow(
		`SELECT *
		 FROM 	(SELECT COUNT(nickname) FROM users) as user_count,
		 		(SELECT COUNT(slug) FROM forum) as forum_count,
				(SELECT COUNT(id) FROM threads) as thread_count,
				(SELECT COUNT(id) FROM posts) as post_count`)
	err := row.Scan(&srvc.User, &srvc.Forum, &srvc.Thread, &srvc.Post)
	if err != nil {
		sr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}
	return srvc, nil
}

func (sr *ServiceRepository) ClearService() error {
	_, err := sr.db.Exec("TRUNCATE users, forum, threads, posts, forum_users, votes;")
	if err != nil {
		sr.logger.Panicln(err.Error())
	}
	return nil
}
