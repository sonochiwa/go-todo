package todo

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	appConfig "go-todo/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

var coll *mongo.Collection
var ctx = context.TODO()
var config = appConfig.GetConfig()

type Todo struct {
	CreatedAt time.Time          `bson:"createdAt"`
	Title     string             `bson:"title"`
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Completed bool               `bson:"completed"`
}

func init() {
	clientOptions := options.Client().ApplyURI(config.MongoDB.MongoDbUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	coll = client.Database(config.MongoDB.DbName).Collection(config.MongoDB.CollectionName)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var results []Todo
	cursor, _ := coll.Find(ctx, bson.D{})
	cursor.All(context.TODO(), &results)
	json.NewEncoder(w).Encode(results)
}

func getTodoById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todo_id")
	idPrimitive, _ := primitive.ObjectIDFromHex(id)
	var t Todo
	coll.FindOne(context.TODO(), bson.D{{"_id", idPrimitive}}).Decode(&t)

	json.NewEncoder(w).Encode(t)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t Todo
	json.NewDecoder(r.Body).Decode(&t)
	document := Todo{Title: t.Title, Completed: false, CreatedAt: time.Now()}
	result, _ := coll.InsertOne(context.TODO(), document)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "todo_id")
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	var t Todo
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

func Routes() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		r.Get("/{todo_id}", getTodoById)
		r.Post("/", createTodo)
		r.Put("/{todo_id}", updateTodo)
		r.Delete("/{todo_id}", deleteTodo)
	})
	return rg
}
