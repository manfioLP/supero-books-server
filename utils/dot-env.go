package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// load .env file
var err = godotenv.Load("./server/.env")


func GetEnvVariable(key string) string {

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}