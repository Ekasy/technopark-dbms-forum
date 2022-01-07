package repository

import (
	"context"
	"database/sql"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/votes"
	"log"
	"regexp"
)

type VoteRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewVoteRepository(db *sql.DB) votes.VoteRepository {
	return &VoteRepository{
		db:     db,
		logger: log.Default(),
	}
}

func (vr *VoteRepository) SelectThread(vote *models.Vote) (int64, error) {
	row := vr.db.QueryRow(
		"SELECT id from threads WHERE 0 = $1 AND slug = $2 OR $2 = '' AND id = $1",
		vote.ThreadId, vote.ThreadSlug)
	err := row.Scan(&vote.ThreadId)
	if err != nil {
		res, _ := regexp.Match(".*no rows in result set.*", []byte(err.Error()))
		if res {
			return 0, myerr.ThreadNotExists
		}
		return 0, myerr.InternalDbError
	}
	return vote.ThreadId, nil
}

func (vr *VoteRepository) UpdateVote(vote *models.Vote) (*models.Thread, error) {
	tx, err := vr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, myerr.InternalDbError
	}

	thread := &models.Thread{}
	// need to detect not existing thread
	res, err := vr.db.Exec(
		"UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3;",
		vote.Voice, vote.Nickname, vote.ThreadId)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}

		res, _ := regexp.Match(".*foreign key constraint \"votes_nickname_fkey\".*", []byte(err.Error()))
		if res {
			return nil, myerr.UserNotExist
		}

		vr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	ra, err := res.RowsAffected()
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}
		vr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	var row *sql.Row
	if ra == 0 {
		_, err = tx.Exec("INSERT INTO votes(voice, nickname, thread) VALUES ($2, $3, $1);", vote.ThreadId, vote.Voice, vote.Nickname)
		if err != nil {
			rollbackError := tx.Rollback()
			if rollbackError != nil {
				return nil, myerr.RollbackError
			}

			res, _ := regexp.Match(".*foreign key constraint \"votes_nickname_fkey\".*", []byte(err.Error()))
			if res {
				return nil, myerr.UserNotExist
			}

			vr.logger.Println(err.Error())
			return nil, myerr.InternalDbError
		}
	}

	row = tx.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1", vote.ThreadId)
	err = row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return nil, myerr.RollbackError
		}
		vr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return nil, myerr.CommitError
	}
	return thread, nil
}
