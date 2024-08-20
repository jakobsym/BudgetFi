package main

import (
	"log"
	"net/http"

	"github.com/jakobsym/BudgetFi/api/internal/controller/budgetfi"
	httphandler "github.com/jakobsym/BudgetFi/api/internal/handler/http"
	"github.com/jakobsym/BudgetFi/api/internal/repository/msmysql"
)

func main() {
	log.Println("starting backend service")
	//repo := mysql.New()
	repo := msmysql.New()
	ctrl := budgetfi.New(repo)
	h := httphandler.New(ctrl)

	//http.HandleFunc("/register", h.CreateUser)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("login/auth", h.OauthCallback) // server side processing route
	http.HandleFunc("/dashboard", h.Login)         // TODO: implement
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
