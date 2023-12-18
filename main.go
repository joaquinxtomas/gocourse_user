package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joaquinxtomas/gocourse_user/internal/user"
	"github.com/joaquinxtomas/gocourse_user/pkg/bootstrap"
	"github.com/joho/godotenv"
)

func main() {

	router := mux.NewRouter()
	_ = godotenv.Load()
	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()

	if err != nil {
		l.Fatal(err)
	}

	userRepo := user.NewRepo(l, db)
	userServ := user.NewService(l, userRepo)
	userEnd := user.MakeEndpoints(userServ)

	router.HandleFunc("/users", userEnd.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEnd.Get).Methods("GET")
	router.HandleFunc("/users", userEnd.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEnd.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEnd.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler:      http.TimeoutHandler(router, time.Second*5, "Tiempo culminado"),
		Addr:         "127.0.0.1:8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err)
	}

}
