package handlers

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Without this check, all routes will yield success response
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Asset not found\n")) // Convert string to byte
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Running API v1\n")) // Convert string to byte
}
