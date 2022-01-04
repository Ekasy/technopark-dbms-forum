package main

import (
	"fmt"
	"forum/db"
	forumdeli "forum/internal/pkg/forum/delivery"
	forumrepo "forum/internal/pkg/forum/repository"
	forumusec "forum/internal/pkg/forum/usecase"
	"forum/internal/pkg/middleware"

	thrddeli "forum/internal/pkg/threads/delivery"
	thrdrepo "forum/internal/pkg/threads/repository"
	thrdusec "forum/internal/pkg/threads/usecase"

	_ "github.com/jackc/pgx/stdlib"

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
	defer db.Close()

	ur := userrepo.NewUserRepository(db)
	uu := userusec.NewUserUsecase(ur)
	ud := userdeli.NewUserDelivery(uu)

	fr := forumrepo.NewForumRepository(db)
	fu := forumusec.NewForumUsecase(fr)
	fd := forumdeli.NewForumDelivery(fu)

	tr := thrdrepo.NewThreadRepository(db)
	tu := thrdusec.NewThreadUsecase(tr)
	td := thrddeli.NewForumDelivery(tu)

	r := mux.NewRouter()
	r = r.PathPrefix("/api").Subrouter()
	r.Use(middleware.ContentTypeMiddleware)
	ud.Routing(r)
	fd.Routing(r)
	td.Routing(r)

	port := 5000
	log.Default().Printf("start serving ::%d\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	log.Default().Fatalf("http serve error %v", err)
}
