package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-todo/pkg/config"
	"go-todo/pkg/todo"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/todos", todo.Routes())

	// TODO: add gracefully shutdown
	http.ListenAndServe(config.Conf.Port, r)
}
