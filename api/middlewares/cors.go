package middlewares

import (
	"github.com/go-chi/cors"
	appConfig "go-todo/configs"
)

var config = appConfig.GetConfig()

func GetCors() cors.Options {
	return cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: config.Cors.AllowCredentials,
		MaxAge:           config.Cors.MaxAge,
	}
}
