package main

import (
	"fmt"
	"forum/db"
	forumdeli "forum/internal/pkg/forum/delivery"
	forumrepo "forum/internal/pkg/forum/repository"
	forumusec "forum/internal/pkg/forum/usecase"
	"forum/internal/pkg/middleware"

	userdeli "forum/internal/pkg/user/delivery"
	userrepo "forum/internal/pkg/user/repository"
	userusec "forum/internal/pkg/user/usecase"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := db.NewDatabase(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", "password", "127.0.0.1", "5432", "forum"))
	if err != nil {
		return
	}

	ur := userrepo.NewUserRepository(db.GetPool())
	uu := userusec.NewUserUsecase(ur)
	ud := userdeli.NewUserDelivery(uu)

	fr := forumrepo.NewForumRepository(db.GetPool())
	fu := forumusec.NewForumUsecase(fr)
	fd := forumdeli.NewForumDelivery(fu)

	r := mux.NewRouter()
	r.Use(middleware.ContentTypeMiddleware)
	ud.Routing(r)
	fd.Routing(r)

	port := 5000
	log.Default().Printf("start serving ::%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	log.Default().Fatalf("http serve error %v", err)
}
