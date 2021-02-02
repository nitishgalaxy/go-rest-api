package main

import (
	"fmt"
	"net/http"
	"os"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Without this check, all routes will yield success response
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Asset not found\n")) // Convert string to byte
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running API v1\n")) // Convert string to byte
}

func main() {
	http.HandleFunc("/", rootHandler)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		// Exit gracefully without showing lot of overwhelming error message
		fmt.Println(err)
		os.Exit(1)
	}
}
