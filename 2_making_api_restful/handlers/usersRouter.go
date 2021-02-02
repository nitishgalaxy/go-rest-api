package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

func UsersRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	path := strings.TrimSuffix(r.URL.Path, "/")

	if path == "/users" {
		switch r.Method {
		case http.MethodGet:
			usersGetAll(w, r)
			return

		case http.MethodPost:
			usersPostOne(w, r)
			return

		default:
			postError(w, http.StatusMethodNotAllowed)
			return
		}
	}

	path = strings.TrimPrefix(r.URL.Path, "/users/")

	if !bson.IsObjectIdHex(path) {
		postError(w, http.StatusNotFound)
		return
	}

	//id := bson.ObjectIdHex(path)

	switch r.Method {
	case http.MethodGet:

		return

	case http.MethodPut:

		return

	case http.MethodPatch:

		return

	case http.MethodDelete:

		return

	default:
		return
	}
}
