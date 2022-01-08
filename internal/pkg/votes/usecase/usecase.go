package usecase

import (
	myerr "forum/internal/error"
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
	err = vu.repo.InsertVote(vote)
	switch err {
	case nil:
	case myerr.ThreadAlreadyExist:
		err = vu.repo.UpdateVote(vote)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	thread, err := vu.repo.SelectThreadById(vote.ThreadId)
	return thread, err
}
