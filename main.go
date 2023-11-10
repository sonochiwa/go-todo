package main

import (
	"context"
	"encoding/json"
	"go-todo/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var coll *mongo.Collection
var ctx = context.TODO()
var conf = config.New()

type todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	Completed bool               `bson:"completed"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func init() {
	clientOptions := options.Client().ApplyURI(conf.MongoDbUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	coll = client.Database(conf.DbName).Collection(conf.CollectionName)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todos", http.StatusSeeOther)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []todo
	cursor, _ := coll.Find(ctx, bson.D{})
	cursor.All(context.TODO(), &results)
	json.NewEncoder(w).Encode(results)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t todo
	json.NewDecoder(r.Body).Decode(&t)
	document := todo{Title: t.Title, Completed: false, CreatedAt: time.Now()}
	result, _ := coll.InsertOne(context.TODO(), document)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todo_id")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	var t todo
	json.NewDecoder(r.Body).Decode(&t)
	update := bson.D{{"$set", bson.D{{"title", t.Title}, {"completed", t.Completed}}}}
	result, err := coll.UpdateOne(ctx, bson.M{"_id": idPrimitive}, update)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(result.ModifiedCount)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todo_id")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	result, err := coll.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(result.DeletedCount)
}

func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		r.Post("/", createTodo)
		r.Put("/{todo_id}", updateTodo)
		r.Delete("/{todo_id}", deleteTodo)
	})
	return rg
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Mount("/todos", todoHandlers())

	// TODO: add gracefully shutdown
	http.ListenAndServe(conf.Port, r)
}
