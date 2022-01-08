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

func (vr *VoteRepository) InsertVote(vote *models.Vote) error {
	tx, err := vr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return myerr.InternalDbError
	}

	_, err = tx.Exec("INSERT INTO votes(voice, nickname, thread) VALUES ($2, $3, $1);", vote.ThreadId, vote.Voice, vote.Nickname)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return myerr.RollbackError
		}

		res, _ := regexp.Match(".*votes_pkey.*", []byte(err.Error()))
		if res {
			return myerr.ThreadAlreadyExist
		}

		res, _ = regexp.Match(".*votes_nickname_fkey.*", []byte(err.Error()))
		if res {
			return myerr.UserNotExist
		}

		vr.logger.Println(err.Error())
		return myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return myerr.CommitError
	}

	return nil
}

func (vr *VoteRepository) UpdateVote(vote *models.Vote) error {
	tx, err := vr.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return myerr.InternalDbError
	}

	_, err = tx.Exec(
		"UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3;",
		vote.Voice, vote.Nickname, vote.ThreadId)
	if err != nil {
		rollbackError := tx.Rollback()
		if rollbackError != nil {
			return myerr.RollbackError
		}

		res, _ := regexp.Match(".*votes_nickname_fkey.*", []byte(err.Error()))
		if res {
			return myerr.UserNotExist
		}

		vr.logger.Println(err.Error())
		return myerr.InternalDbError
	}

	err = tx.Commit()
	if err != nil {
		return myerr.CommitError
	}

	return nil
}

func (vr *VoteRepository) SelectThreadById(threadId int64) (*models.Thread, error) {
	thread := &models.Thread{}

	row := vr.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1", threadId)
	err := row.Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	if err != nil {
		vr.logger.Println(err.Error())
		return nil, myerr.InternalDbError
	}

	return thread, nil
}
