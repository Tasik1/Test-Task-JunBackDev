package main

import (
	"TestBackDev/route"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	loadenv()
	log.Fatal(route.RunAPI(":8080"))
}

func loadenv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading .env file: " + err.Error())
	}
}
