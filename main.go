package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	appConfig "go-todo/pkg/config"
	"go-todo/pkg/todo"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var config = appConfig.GetConfig()

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(corsMiddleware.Handler)
	r.Mount("/todos", todo.Routes())

	srv := &http.Server{
		Addr:         config.Port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port", config.Port)
		if err := http.ListenAndServe(config.Port, r); err != nil {
			log.Printf("listen:%s\n", err)
		}
	}()

	<-ch
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped")
}
