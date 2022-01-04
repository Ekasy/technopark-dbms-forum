package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/forum"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type ForumDelivery struct {
	forumUsecase forum.ForumUsecase
}

func NewForumDelivery(forumUsecase forum.ForumUsecase) *ForumDelivery {
	return &ForumDelivery{
		forumUsecase: forumUsecase,
	}
}

func (fd *ForumDelivery) Routing(r *mux.Router) {
	r.HandleFunc("/forum/create", fd.CreateForumHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/forum/{slug}/details", fd.GetForumHandler).Methods(http.MethodGet, http.MethodOptions)
}

func (fd *ForumDelivery) CreateForumHandler(w http.ResponseWriter, r *http.Request) {
	forumInput := &models.ForumInput{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, forumInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	forum := forumInput.ToDefaultForum()
	forum, err = fd.forumUsecase.CreateForum(forum)
	switch err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		w.Write(models.ToBytes(forum))
	case myerr.UserNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("Can't find forum's owner: %s", forum.User)}))
	case myerr.ForumAlreadyExist:
		w.WriteHeader(http.StatusConflict)
		w.Write(models.ToBytes(forum))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (fd *ForumDelivery) GetForumHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	res, _ := regexp.Match("[a-z0-9]+(?:-[a-z0-9]+)*", []byte(slug))
	if !res {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid slug"}))
		return
	}

	forum, err := fd.forumUsecase.GetForum(slug)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(forum))
	case myerr.NoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("forum %s not found", slug)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
