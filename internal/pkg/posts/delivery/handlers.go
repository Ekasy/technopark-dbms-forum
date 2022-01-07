package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/posts"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type PostDelivery struct {
	postUsecase posts.PostUsecase
}

func NewPostDelivery(postUsecase posts.PostUsecase) *PostDelivery {
	return &PostDelivery{
		postUsecase: postUsecase,
	}
}

func (pd *PostDelivery) Routing(r *mux.Router) {
	r.HandleFunc("/thread/{slug_or_id}/create", pd.CreatePostHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/thread/{slug_or_id}/posts", pd.GetPostsByThreadHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/post/{id}/details", pd.GetPostDetailHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/post/{id}/details", pd.UpdatePostHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (pd *PostDelivery) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug_or_id"]
	id, err := strconv.ParseInt(slug, 10, 64)
	if err == nil {
		slug = ""
	} else {
		id = 0
	}
	postsInput := []*models.PostInput{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, &postsInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	posts, err := pd.postUsecase.CreatePostsBySlugOrId(slug, id, postsInput)
	switch err {
	case nil:
		w.WriteHeader(http.StatusCreated)
		w.Write(models.ToBytes(posts))
	case myerr.ThreadNotExists:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("thread {slug: %s, id: %d} not found", slug, id)}))
	case myerr.ParentNotExist:
		w.WriteHeader(http.StatusConflict)
		w.Write(models.ToBytes(models.Error{Message: "one parent not found"}))
	case myerr.UserNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: "one user not found"}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (pd *PostDelivery) GetPostsByThreadHandler(w http.ResponseWriter, r *http.Request) {
	tq := models.NewThreadQuery(mux.Vars(r), r.URL.Query())
	posts, err := pd.postUsecase.GetPostsRec(tq)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(posts))
	case myerr.ThreadNotExists:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("thread {slug: '%s', id: %d} not found", tq.ThreadSlug, tq.ThreadId)}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (pd *PostDelivery) GetPostDetailHandler(w http.ResponseWriter, r *http.Request) {
	pq := models.NewPostQuery(mux.Vars(r), r)
	info, err := pd.postUsecase.GetInfo(pq)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(info))
	case myerr.ThreadNotExists:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: "thread not found"}))
	case myerr.ForumNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: "forum not found"}))
	case myerr.PostNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: "post not found"}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (pd *PostDelivery) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	pu := &models.PostUpdate{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, &pu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err == nil {
		pu.Id = id
	}
	post, err := pd.postUsecase.UpdatePost(pu)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(post))
	case myerr.PostNotExist:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: "post not found"}))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
