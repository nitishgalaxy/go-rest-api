package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nitishgalaxy/go-rest-api/handlers"
)

func main() {
	http.HandleFunc("/users", handlers.UsersRouter)
	http.HandleFunc("/users/", handlers.UsersRouter)
	http.HandleFunc("/", handlers.RootHandler)

	fmt.Println("Server started at localhost:8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		// Exit gracefully without showing lot of overwhelming error message
		fmt.Println(err)
		os.Exit(1)
	}
}
