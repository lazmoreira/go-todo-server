package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lazmoreira/go-todo/router"
)

func main() {
	r := router.Router()

	fmt.Println("Starting on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
