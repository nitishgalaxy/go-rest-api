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

		case http.MethodHead:
			usersGetAll(w, r)
			return

		case http.MethodOptions:
			postOptionsResponse(w, []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}, nil)
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

	id := bson.ObjectIdHex(path)

	switch r.Method {
	case http.MethodGet:
		usersGetOne(w, r, id)
		return

	case http.MethodHead:
		usersGetOne(w, r, id)
		return

	case http.MethodPut:
		usersPutOne(w, r, id)
		return

	case http.MethodPatch:
		usersPatchOne(w, r, id)
		return

	case http.MethodDelete:
		usersDeleteOne(w, r, id)
		return

	case http.MethodOptions:
		postOptionsResponse(w, []string{http.MethodGet, http.MethodPut, http.MethodHead, http.MethodOptions}, nil)
		return

	default:
		return
	}
}
