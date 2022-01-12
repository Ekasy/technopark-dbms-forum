package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/forum"
	"io/ioutil"
	"net/http"

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
	r.HandleFunc("/forum/{slug}/users", fd.GetUsersHandler).Methods(http.MethodGet, http.MethodOptions)
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

func (fd *ForumDelivery) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	fv := models.NewForumUsersQuery(mux.Vars(r), r.URL.Query())
	users, err := fd.forumUsecase.GetUsersByForum(fv)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(users))
	case myerr.ForumNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("forum %s not found", fv.ForumSlug)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
