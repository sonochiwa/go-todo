package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-todo/pkg/config"
	"go-todo/pkg/todo"
	"net/http"
)

func main() {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешенные источники
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Максимальное время кэширования предопределенных ответов в секундах.
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(corsMiddleware.Handler)

	r.Mount("/todos", todo.Routes())

	// TODO: add gracefully shutdown
	http.ListenAndServe(config.Conf.Port, r)
}
