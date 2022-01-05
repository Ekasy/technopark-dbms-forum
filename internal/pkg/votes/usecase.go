package votes

import "forum/internal/models"

type VoteUsecase interface {
	UpdateVote(vote *models.Vote) (*models.Thread, error)
}
