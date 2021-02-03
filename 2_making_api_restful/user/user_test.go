package user

import (
	"os"
	"reflect"
	"testing"

	"github.com/asdine/storm"
	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	m.Run() // Execute all tests for given package
	os.Remove(dbpath)
}

func TestCrud(t *testing.T) {
	// Log the name of the test
	// t.Log function ensures only failing tests display this info.
	t.Log("Create")
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}

	err := u.Save()
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error saving a record: %s ", err)
	}

	//=========================================

	t.Log("Read")
	u2, err := One(u.ID)
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error receiving a record: %s ", err)
	}

	if !reflect.DeepEqual(u2, u) {
		t.Error("Records do not match")
	}

	//=========================================

	t.Log("Update")
	u.Role = "developer"
	err = u.Save()
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error updating a record: %s ", err)
	}

	u3, err := One(u.ID)
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error receiving a record: %s ", err)
	}

	if !reflect.DeepEqual(u3, u) {
		t.Error("Records do not match")
	}

	//=========================================

	t.Log("Delete")
	err = Delete(u.ID)
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error deleting a record: %s ", err)
	}

	_, err = One(u.ID)
	if err == nil {
		t.Fatal("Error should not exist anymore")
	}

	if err != storm.ErrNotFound {
		t.Fatalf("Error retrieving non existing record: %s", err)
	}

	//=========================================

	t.Log("Read All")
	u.ID = bson.NewObjectId()
	u2.ID = bson.NewObjectId()
	u3.ID = bson.NewObjectId()

	err = u.Save()
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error saving a record: %s ", err)
	}

	err = u2.Save()
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error saving a record: %s ", err)
	}

	err = u3.Save()
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error saving a record: %s ", err)
	}

	users, err := All()
	if err != nil {
		// t.Fatalf/ t.Fatal is shortcut for t.FatalNow
		t.Fatalf("Error saving all records: %s ", err)
	}

	if len(users) != 3 {
		t.Errorf("Different no. of records retrieved. Expected 3 / Actual %d", len(users))
	}
}
