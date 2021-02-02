package main

import "net/http"

func main() {
	err := http.ListenAndServe("localhost:11111", nil)
	if err != nil {
		panic(err) // Terminate the program and display error stack trace
	}
}
