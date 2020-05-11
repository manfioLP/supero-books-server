package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

func Pagination(r *http.Request) (int, int) {

	page, err := strconv.Atoi(r.FormValue("page"))
	if err == nil {
		fmt.Println("error casting page queryString to number, USING DEFAULT VALUE")
		page = 10
	}

	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err == nil {
		fmt.Println("error casting limit queryString to number, USING DEFAULT VALUE")
		limit = 10
	}

	//if page == 0 {
	//	page = 1
	//}

	if limit == 0 {
		limit = 10
	}

	return page, limit
}

func FindFilter(r *http.Request) {
	filter  := bson.M{}

	title := r.FormValue("title")
	author := r.FormValue("author")
	isbn := r.FormValue("isbn")

	if title != ""{
		filter = bson.M{title: title}
	} else if author != ""{
		filter = bson.M{author: author}
	} else if isbn != ""{
		filter = bson.M{isbn: isbn}
	}

	return filter
}
