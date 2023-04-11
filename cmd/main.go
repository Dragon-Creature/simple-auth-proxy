package main

import (
	"git.ssns.se/git/frozendragon/simple-auth-proxy/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/*", auth.GetLoginPage)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
