package delivery

import (
	"encoding/json"
	"fmt"
	myerr "forum/internal/error"
	"forum/internal/models"
	"forum/internal/pkg/user"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type UserDelivery struct {
	userUsecase user.UserUsecase
}

func NewUserDelivery(userUsecase user.UserUsecase) *UserDelivery {
	return &UserDelivery{
		userUsecase: userUsecase,
	}
}

func (ud *UserDelivery) Routing(r *mux.Router) {
	r.HandleFunc("/user/{nickname}/create", ud.CreateUserHandler).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/user/{nickname}/profile", ud.GetUserHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/user/{nickname}/profile", ud.UpdateUserHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (ud *UserDelivery) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	userInput := &models.UserUpdate{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	users, inserted, err := ud.userUsecase.CreateUser(userInput.ToUser(nickname))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
		return
	}

	if inserted {
		w.WriteHeader(http.StatusCreated)
		w.Write(models.ToBytes(users[0]))
	} else {
		w.WriteHeader(http.StatusConflict)
		w.Write(models.ToBytes(users))
	}
}

func (ud *UserDelivery) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	user, err := ud.userUsecase.GetUser(nickname)
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(user))
	case myerr.NoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("Can't find user with nickname %s", nickname)}))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (ud *UserDelivery) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	nickname := mux.Vars(r)["nickname"]
	userInput := &models.UserUpdate{}
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 1"}))
		return
	}

	err = json.Unmarshal(buf, userInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: "invalid body 2"}))
		return
	}

	user, err := ud.userUsecase.UpdateUser(userInput.ToUser(nickname))
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(user))
	case myerr.NoRows:
		w.WriteHeader(http.StatusNotFound)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("Can't find user with nickname %s", nickname)}))
	case myerr.EmailAlreadyExist:
		w.WriteHeader(http.StatusConflict)
		w.Write(models.ToBytes(models.Error{Message: fmt.Sprintf("Can't update email for user with nickname %s", nickname)}))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
