package user

import (
	"errors"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

// User holds data for single user
type User struct {
	ID   bson.ObjectId `json:"id" storm:"id"`
	Name string        `json:"name"`
	Role string        `json:"role"`
}

const (
	dbpath = "users.db"
)

var (
	ErrRecordInvalid = errors.New("User record is invalid")
)

// All() retrrieves all users from Database
func All() ([]User, error) {
	// Open DB Connection
	db, err := storm.Open(dbpath)

	if err != nil {
		return nil, err
	}

	defer db.Close()
	users := []User{}
	err = db.All(&users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// One() retrrieves single user record from Database
func One(id bson.ObjectId) (*User, error) {
	// Open DB Connection
	db, err := storm.Open(dbpath)

	if err != nil {
		return nil, err
	}

	defer db.Close()
	user := new(User)            // Returns pointer
	err = db.One("ID", id, user) // user is already a pointer

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Delete() removes single user record from Database
func Delete(id bson.ObjectId) error {
	// Open DB Connection
	db, err := storm.Open(dbpath)

	if err != nil {
		return err
	}

	defer db.Close()

	// In order to delete an object, storm requires us to retrieve it first
	user := new(User)            // Returns pointer
	err = db.One("ID", id, user) // user is already a pointer

	if err != nil {
		return err
	}

	return db.DeleteStruct(user)
}

// Save() updates or creates a given record in the Database
func (user *User) Save() error {
	// Shorthand sytax
	if err := user.validate(); err != nil {
		return err
	}

	// Open DB Connection
	db, err := storm.Open(dbpath)

	if err != nil {
		return err
	}

	defer db.Close()

	return db.Save(user)
}

// Validation logic to make sure function contains valid data
func (user *User) validate() error {
	if user.Name == "" {
		return ErrRecordInvalid
	}

	return nil // No error
}
