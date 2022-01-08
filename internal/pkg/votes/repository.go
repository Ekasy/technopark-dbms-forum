package votes

import "forum/internal/models"

type VoteRepository interface {
	InsertVote(vote *models.Vote) error
	UpdateVote(vote *models.Vote) error
	SelectThreadById(threadId int64) (*models.Thread, error)
	SelectThread(vote *models.Vote) (int64, error)
}
