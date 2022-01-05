package votes

import "forum/internal/models"

type VoteRepository interface {
	UpdateVote(vote *models.Vote) (*models.Thread, error)
	SelectThread(vote *models.Vote) (int64, error)
}
