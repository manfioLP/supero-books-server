package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"../models"
	"../utils"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	. "github.com/gobeam/mongo-go-pagination"
)

// DB connection string, if local use localhost below
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
	fmt.Println(connectionString)
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
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	page, limit := utils.Pagination(r)
	filter := utils.Filter(r)

	payload, paginationInfo := getAllBooks(limit, page, filter)
	fmt.Println(paginationInfo)
	answer := models.GetAnswer{
		Ok:     true,
		Data:   payload,
		Page: paginationInfo.Page,
		Total: paginationInfo.Total,
		PerPage: paginationInfo.PerPage,

	}
	json.NewEncoder(w).Encode(answer)
}

func RegisterBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var book models.Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	v := validator.New()
	err := v.Struct(book)

	var errs []string
	var answer models.Answer

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, utils.HandleDbError(e))
		}
		fmt.Println("validated...")
		answer = models.Answer{
			Ok:     false,
			Errors: errs,
		}
		fmt.Println(errs)
		w.WriteHeader(400)
	} else {
		book.ID = book.ISBN
		if _, insertError := insertBook(book); insertError != "" {
			w.WriteHeader(409)
			answer = models.Answer{
				Ok: false,
				Errors: append(errs, insertError),
			}
		} else {
			answer = models.Answer{
				Ok:     true,
				Errors: nil,
				Data:   &book,
			}
			fmt.Println(book, r.Body)
		}

	}

	json.NewEncoder(w).Encode(answer)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)
	book := getBook(params["ISBN"])
	json.NewEncoder(w).Encode(book)
}

func DeleteBook (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	deleteBook(params["ISBN"])

	json.NewEncoder(w).Encode(models.Answer{ Ok: true })
}

func getAllBooks(limit int, page int, filter bson.M) ([]models.Book, PaginationData) {

	// TODO: add sorting by year

	projection := bson.D{
		{"title", 1},
		{"author", 1},
		{"isbn", 1},
		{"publisher", 1},
		{"year", 1},
		{"language", 1},
	}

	// sorted data use below
	// paginatedData, err := New(collection).Limit(limit).Page(page).Sort("year", -1).Select(projection).Filter(filter).Find()
	paginatedData, err := New(collection).Limit(int64(limit)).Page(int64(page)).Select(projection).Filter(filter).Find()
	if err != nil {
		fmt.Println("Error on paginated data /getallbooks")
		log.Fatal(err)
	}
	var lists []models.Book
	for _, raw := range paginatedData.Data {
		var book *models.Book
		if marshallErr := bson.Unmarshal(raw, &book); marshallErr == nil {
			lists = append(lists, *book)
		}
	}

	return lists, paginatedData.Pagination
}

func getBook(ISBN string) models.Book {
	var book models.Book
	filter := bson.M{"_id": ISBN}
	if err := collection.FindOne(context.Background(), filter).Decode(&book); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Single book ISBN#" + ISBN + " found:", book)
	return book
}

func insertBook(book models.Book) (interface{}, string) {
	insertResult, err := collection.InsertOne(context.Background(), book)

	var errmsg = ""
	var insertedId interface{} = nil

	if err != nil {
		fmt.Println(err)
		errmsg = utils.HandleWriteError(err)
	} else {
		insertedId = insertResult.InsertedID
	}

	return insertedId, errmsg
}

func deleteBook(ISBN string) {
	filter := bson.M{"_id": ISBN}
	d, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted Document", d.DeletedCount)
}