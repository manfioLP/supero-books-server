package main

import (
	"./router"
	"./utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := utils.GetEnvVariable("PORT")
	r:= router.Router()

	fmt.Println("Starting Server on port " + port + "....")

	log.Fatal(http.ListenAndServe(":" + port, r))
}
