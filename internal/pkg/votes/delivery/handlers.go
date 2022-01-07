package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/votes"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type VoteDelivery struct {
	voteUsecase votes.VoteUsecase
}

func NewVoteDelivery(voteUsecase votes.VoteUsecase) *VoteDelivery {
	return &VoteDelivery{
		voteUsecase: voteUsecase,
	}
}

func (vd *VoteDelivery) Routing(r *mux.Router) {
	r.HandleFunc("/thread/{slug_or_id}/vote", vd.UpdateVoteHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (vd *VoteDelivery) UpdateVoteHandler(w http.ResponseWriter, r *http.Request) {
	vote := &models.Vote{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, &vote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	vote.ThreadSlug = mux.Vars(r)["slug_or_id"]
	vote.ThreadId, err = strconv.ParseInt(vote.ThreadSlug, 10, 64)
	if err != nil {
		vote.ThreadId = 0
	} else {
		vote.ThreadSlug = ""
	}

	thread, err := vd.voteUsecase.UpdateVote(vote)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(thread))
	case myerr.ThreadNotExists:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("thread {slug: %s, id: %d} not found", vote.ThreadSlug, vote.ThreadId)}))
	case myerr.UserNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("user %s not found", vote.Nickname)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
