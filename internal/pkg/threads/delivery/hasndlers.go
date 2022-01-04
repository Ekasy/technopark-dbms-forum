package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/threads"
	"io/ioutil"
	"net/http"

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
