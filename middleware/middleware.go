package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"../models"
	"../utils"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/ DB connection string, if local use localhost below
// const connectionString = "mongodb://localhost:27017"
var dbPwd = utils.GetEnvVariable("DB_PASSWORD")
var dbUsername = utils.GetEnvVariable("DB_USERNAME")

var connectionString = "mongodb+srv://" + dbUsername + ":" + dbPwd + "@cluster0-jtohn.mongodb.net/test?retryWrites=true&w=majority"

const dbName = "supero-books"

// Collection's name
const collName = "books"

// instance of collection obj
var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection = client.Database(dbName).Collection(collName)

	fmt.Println("Collection instance created!")
}

// DB Functions - CRUDDING

// get all books
// todo: add params
// todo: pagination
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllBooks()
	json.NewEncoder(w).Encode(payload)
}

func RegisterBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	v := validator.New()
	err := v.Struct(book)

	// todo: break condition if err
	fmt.Println("validation-err: ", err)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e)
	}

	fmt.Println(book, r.Body)
	insertBook(book)
	json.NewEncoder(w).Encode(book)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// todo: check for isbn or id
	params := mux.Vars(r)
	//book := getBook(params["ID"])
	book := getBook(params["ISBN"])
	json.NewEncoder(w).Encode(book)
}

func DeleteBook (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	//deleteBook(params["ID"])
	deleteBook(params["ISBN"])
	json.NewEncoder(w).Encode(params["id"])
}

// todo: add pagination
func getAllBooks() []primitive.M {

	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	fmt.Println("cur books... :")
	fmt.Println(cur)
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results

}

func getBook(ISBN string) models.Book {
	var book models.Book
	filter := bson.M{"ISBN": ISBN}
	if err := collection.FindOne(context.Background(), filter).Decode(&book); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Single book ISBN#" + ISBN + " found:", book)
	return book
}

func insertBook(book models.Book) {
	insertResult, err := collection.InsertOne(context.Background(), book)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted a single record", insertResult)
}

func deleteBook(ISBN string) {
	fmt.Println(ISBN)
	//id, _ := primitive.ObjectIDFromHex(task)
	//filter := bson.M{"_id": id}
	filter := bson.M{"ISBN": ISBN}
	d, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
}