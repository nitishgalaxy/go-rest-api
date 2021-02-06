package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/nitishgalaxy/go-rest-api-echo/cache"
	"github.com/nitishgalaxy/go-rest-api-echo/user"
	"gopkg.in/mgo.v2/bson"
)

type jsonResponse map[string]interface{}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Running API v1")
}

func usersOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func userOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodHead, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func bodyToUser(r *http.Request, u *user.User) error {
	// Handle edge cases
	if r == nil {
		return errors.New("a request is required")
	}

	if r.Body == nil {
		return errors.New("request body is empty")
	}

	if u == nil {
		return errors.New("a user is required")
	}

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bd, u)
}

func usersPutOne(c echo.Context) error {
	u := new(user.User)
	//err := bodyToUser(r, u)
	err := c.Bind(u)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	u.ID = id
	err = u.Save()

	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	cache.Drop("/users")
	//cache.Drop(cache.MakeResource(r))
	//cw := cache.NewWriter(w, r)
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPatchOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	// Get existing record
	u, err := user.One(id)

	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	//err = bodyToUser(r, u)
	err = c.Bind(u)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	u.ID = id
	err = u.Save()

	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

	}

	cache.Drop("/users")
	//cache.Drop(cache.MakeResource(r))
	//cw := cache.NewWriter(w, r)
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPostOne(c echo.Context) error {
	u := new(user.User)
	//err := bodyToUser(r, u)
	err := c.Bind(u)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	u.ID = bson.NewObjectId()
	err = u.Save()

	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	cache.Drop("/users")
	c.Response().Header().Set("Location", "/users/"+u.ID.Hex())
	return c.NoContent(http.StatusCreated)
}

func usersGetOne(c echo.Context) error {
	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}

	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	u, err := user.One(id)

	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError)

	}

	if c.Request().Method == http.MethodHead {
		return c.JSON(http.StatusOK, jsonResponse{}) // Send empty content

	}

	// cw = Cache Writer
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersGetAll(c echo.Context) error {

	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}

	users, err := user.All()

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}

	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"users": users})
}

func usersDeleteOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	id := bson.ObjectIdHex(c.Param("id"))

	// Get existing record
	err := user.Delete(id)

	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	cache.Drop("/users")
	cache.Drop(cache.MakeResource(c.Request()))
	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()

	e.GET("/", root)

	// Echo allows us to group routes and save us typing
	u := e.Group("/users")
	u.OPTIONS("", usersOptions)
	u.HEAD("", usersGetAll)
	u.GET("", usersGetAll)
	u.POST("", usersPostOne)

	// To access individual items from users collection,
	// we will create a subgroup 'uid' from Group u
	uid := u.Group("/:id")
	uid.OPTIONS("", userOptions)
	uid.HEAD("", usersGetOne)
	uid.GET("", usersGetOne)
	uid.PUT("", usersPostOne)
	uid.PATCH("", usersPatchOne)
	uid.DELETE("", usersDeleteOne)

	e.Start(":8080")
}