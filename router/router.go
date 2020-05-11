package router

import (
	"../middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/book", middleware.GetAllBooks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/book", middleware.GetAllBooks).Methods("GET", "OPTIONS").Queries("author", "title", "isbn", "page", "limit")
	router.HandleFunc("/api/book", middleware.RegisterBook).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/book/{ISBN}", middleware.GetBook).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/book/{ISBN}", middleware.DeleteBook).Methods("DELETE", "OPTIONS")
	return router
}