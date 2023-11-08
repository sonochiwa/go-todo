package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// перемещу константы в переменные окружения
const (
	hostName       string = "localhost:27017"
	dbName         string = "go_todo_list"
	collectionName string = "go_todo_list"
	port           string = ":9000"
)

type (
	todoModel struct {
		ID        string    `bson:"_id,omitempty"`
		Title     string    `bson:"title"`
		Completed bool      `bson:"completed"`
		CreatedAt time.Time `bson:"createAt"`
	}

	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("fetchTodos")

	//todos := []todoModel{}
	//
	//if err := db.C(collectionName).
	//	Find(bson.M{}).
	//	All(&todos); err != nil {
	//	rnd.JSON(w, http.StatusProcessing, renderer.M{
	//		"message": "Failed to fetch todo",
	//		"error":   err,
	//	})
	//	return
	//}
	//
	//todoList := []todo{}
	//for _, t := range todos {
	//	todoList = append(todoList, todo{
	//		ID:        t.ID.Hex(),
	//		Title:     t.Title,
	//		Completed: t.Completed,
	//		CreatedAt: t.CreatedAt,
	//	})
	//}
	//
	//rnd.JSON(w, http.StatusOK, renderer.M{
	//	"data": todoList,
	//})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createTodo")

	//var t todo
	//
	//if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
	//	rnd.JSON(w, http.StatusProcessing, err)
	//	return
	//}
	//
	//if t.Title == "" {
	//	rnd.JSON(w, http.StatusBadRequest, renderer.M{
	//		"message": "The title field is required",
	//	})
	//	return
	//}
	//
	//tm := todoModel{
	//	ID:        bson.NewObjectId(),
	//	Title:     t.Title,
	//	Completed: false,
	//	CreatedAt: time.Now(),
	//}
	//if err := db.C(collectionName).Insert(&tm); err != nil {
	//	rnd.JSON(w, http.StatusProcessing, renderer.M{
	//		"message": "Failed to save todo",
	//		"error":   err,
	//	})
	//	return
	//}
	//
	//rnd.JSON(w, http.StatusCreated, renderer.M{
	//	"message": "Todo created successfully",
	//	"todo_id": tm.ID.Hex(),
	//})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateTodo")

	//id := strings.TrimSpace(chi.URLParam(r, "id"))
	//
	//if !bson.IsObjectIdHex(id) {
	//	rnd.JSON(w, http.StatusBadRequest, renderer.M{
	//		"message": "The id is invalid",
	//	})
	//	return
	//}
	//
	//var t todo
	//
	//if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
	//	rnd.JSON(w, http.StatusProcessing, err)
	//	return
	//}
	//
	//if t.Title == "" {
	//	rnd.JSON(w, http.StatusBadRequest, renderer.M{
	//		"message": "The title field is required",
	//	})
	//	return
	//}
	//
	//if err := db.C(collectionName).
	//	Update(
	//		bson.M{"_id": bson.ObjectIdHex(id)},
	//		bson.M{"title": t.Title, "completed": t.Completed},
	//	); err != nil {
	//	rnd.JSON(w, http.StatusProcessing, renderer.M{
	//		"message": "Failed to update todo",
	//		"error":   err,
	//	})
	//	return
	//}
	//
	//rnd.JSON(w, http.StatusOK, renderer.M{
	//	"message": "Todo updated successfully",
	//})
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleteTodo")

	//id := strings.TrimSpace(chi.URLParam(r, "id"))
	//
	//if !bson.IsObjectIdHex(id) {
	//	rnd.JSON(w, http.StatusBadRequest, renderer.M{
	//		"message": "The id is invalid",
	//	})
	//	return
	//}
	//
	//if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err != nil {
	//	rnd.JSON(w, http.StatusProcessing, renderer.M{
	//		"message": "Failed to delete todo",
	//		"error":   err,
	//	})
	//	return
	//}
	//
	//rnd.JSON(w, http.StatusOK, renderer.M{
	//	"message": "Todo deleted successfully",
	//})
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", fetchTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func main() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Listening on port ", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-ch
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}
