package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	appConfig "go-todo/pkg/config"
	"go-todo/pkg/middlewares"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.New(middlewares.InitCors()).Handler)
	r.Mount("/todos", todo.Routes())

	srv := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.IdleTimeout) * time.Second,
	}

	go func() {
		log.Println("Already available on", "http://"+srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	<-ch
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server was stopped")
}
