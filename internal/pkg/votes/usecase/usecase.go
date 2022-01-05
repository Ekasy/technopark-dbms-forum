package usecase

import (
	"forum/internal/models"
	"forum/internal/pkg/votes"
)

type VoteUsecase struct {
	repo votes.VoteRepository
}

func NewVoteUsecase(repo votes.VoteRepository) votes.VoteUsecase {
	return &VoteUsecase{
		repo: repo,
	}
}

func (vu *VoteUsecase) UpdateVote(vote *models.Vote) (*models.Thread, error) {
	_, err := vu.repo.SelectThread(vote)
	if err != nil {
		return nil, err
	}
	thread, err := vu.repo.UpdateVote(vote)
	return thread, err
}
