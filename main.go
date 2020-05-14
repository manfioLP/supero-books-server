package main

import (
	"./router"
	"./utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	port := utils.GetEnvVariable("PORT")
	r:= router.Router()

	fmt.Println("Starting Server on port " + port + "....")

	log.Fatal(http.ListenAndServe(":" + port, handlers.CORS()(r)))
}
