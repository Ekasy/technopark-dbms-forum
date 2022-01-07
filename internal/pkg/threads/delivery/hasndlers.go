package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/threads"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ThreadDelivery struct {
	threadUsecase threads.ThreadUsecase
}

func NewForumDelivery(threadUsecase threads.ThreadUsecase) *ThreadDelivery {
	return &ThreadDelivery{
		threadUsecase: threadUsecase,
	}
}

func (td *ThreadDelivery) Routing(r *mux.Router) {
	r.HandleFunc("/forum/{slug}/create", td.CreateThreadHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/forum/{slug}/threads", td.GetThreadsHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/thread/{slug_or_id}/details", td.GetThreadHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/thread/{slug_or_id}/details", td.UpdateThreadHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (td *ThreadDelivery) CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	thredInput := &models.ThreadInput{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, thredInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	thread, err := td.threadUsecase.CreateThread(thredInput.ToThread(slug))
	switch err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		w.Write(models.ToBytes(thread))
	case myerr.AuthorNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("user %s not found", thredInput.Author)}))
	case myerr.ForumNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("forum %s not found", thredInput.Forum)}))
	case myerr.ThreadAlreadyExist:
		w.WriteHeader(http.StatusConflict)
		w.Write(models.ToBytes(thread))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (td *ThreadDelivery) GetThreadsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tv := models.NewThreadsVars(mux.Vars(r), query)

	threads, err := td.threadUsecase.GetThreadsByForum(tv)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(threads))
	case myerr.NoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("forum %s not exist", tv.ForumSlug)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (td *ThreadDelivery) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	tv := models.NewThreadsVars(mux.Vars(r), query)

	threads, err := td.threadUsecase.GetUsersByForum(tv)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(threads))
	case myerr.NoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("forum %s not exist", tv.ForumSlug)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (td *ThreadDelivery) GetThreadHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug_or_id"]
	id, err := strconv.ParseInt(slug, 10, 64)
	if err == nil {
		slug = ""
	} else {
		id = 0
	}

	thread, err := td.threadUsecase.GetThread(slug, id)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(thread))
	case myerr.ThreadNotExists:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("thread with {id: %d, slug: '%s'} not exist", id, slug)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (td *ThreadDelivery) UpdateThreadHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug_or_id"]
	id, err := strconv.ParseInt(slug, 10, 64)
	if err == nil {
		slug = ""
	} else {
		id = 0
	}

	thredUpdate := &models.ThreadUpdate{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, &thredUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	thredUpdate.Id = id
	thredUpdate.Slug = slug
	thread, err := td.threadUsecase.UpdateThread(thredUpdate)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(thread))
	case myerr.ThreadNotExists:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("thread with {id: %d, slug: '%s'} not exist", id, slug)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
