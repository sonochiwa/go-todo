package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-todo/api/handlers"
	"go-todo/api/middlewares"
	appConfig "go-todo/configs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var config = appConfig.GetConfig()

func mainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.New(middlewares.GetCors()).Handler)
	r.Get("/", mainPage)
	r.Mount("/todos", todo_handlers.Routes())

	srv := &http.Server{
		Addr:         config.Server.Host + ":" + config.Server.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.IdleTimeout) * time.Second,
	}

	go func() {
		log.Println("Already available")
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
