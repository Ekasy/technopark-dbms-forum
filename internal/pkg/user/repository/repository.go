package repository

import (
	"context"
	"database/sql"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/user"
	"log"
	"regexp"
)

type UserRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &UserRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (ur *UserRepository) InsertUser(user *models.User) error {
	tx, err := ur.db.BeginTx(context.Background(), nil)
	if err != nil {
		return myerr.InternalDbError
	}

	row := tx.QueryRowContext(
		context.Background(),
		"INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4) RETURNING nickname, fullname, about, email;",
		user.Nickname, user.Fullname, user.About, user.Email,
	)

	err = row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		res, _ := regexp.Match(".*users_pkey.*", []byte(err.Error()))
		if res {
			return myerr.NicknameAlreadyExist
		}

		res, _ = regexp.Match(".*users_email_key.*", []byte(err.Error()))
		if res {
			return myerr.EmailAlreadyExist
		}
	}

	if err != nil {
		ur.logger.Printf(err.Error())
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

func (ur *UserRepository) UpdateUser(user *models.User) error {
	tx, err := ur.db.BeginTx(context.Background(), nil)
	if err != nil {
		return myerr.InternalDbError
	}

	row := tx.QueryRowContext(
		context.Background(),
		"UPDATE users SET fullname = $2, about = $3, email = $4 WHERE nickname = $1 RETURNING nickname, fullname, about, email;",
		user.Nickname, user.Fullname, user.About, user.Email,
	)

	err = row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return myerr.NoRows
		}

		res, _ = regexp.Match(".*users_email_key.*", []byte(err.Error()))
		if res {
			return myerr.EmailAlreadyExist
		}
	}

	if err != nil {
		ur.logger.Printf(err.Error())
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

func (ur *UserRepository) SelectUser(nickname string) (*models.User, error) {
	row := ur.db.QueryRowContext(
		context.Background(),
		"SELECT nickname, fullname, about, email FROM users WHERE nickname = $1",
		nickname,
	)

	user := &models.User{}
	err := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return nil, myerr.NoRows
		}
	}

	if err != nil {
		ur.logger.Printf(err.Error())
		return nil, myerr.InternalDbError
	}

	return user, nil
}

func (ur *UserRepository) SelectUsersIfExists(nickname string, email string) ([]*models.User, error) {
	rows, err := ur.db.QueryContext(
		context.Background(),
		"SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2",
		nickname, email,
	)
	if err != nil {
		ur.logger.Printf(err.Error())
		return nil, myerr.InternalDbError
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			ur.logger.Printf(err.Error())
			return nil, myerr.InternalDbError
		}
		users = append(users, user)
	}

	return users, nil
}
