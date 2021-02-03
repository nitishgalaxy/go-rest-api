package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/nitishgalaxy/go-rest-api/user"
	"gopkg.in/mgo.v2/bson"
)

func TestBodyToUser(t *testing.T) {
	// Create test data
	valid := &user.User{
		ID:   bson.NewObjectId(),
		Name: "John",
		Role: "Tester",
	}

	js, err := json.Marshal(valid)
	if err != nil {
		t.Errorf("Error marshaalling a valid user: %s", err)
		t.FailNow()
	}

	// ts = Test Suite
	// A list of struct
	ts := []struct {
		txt string // Describe the test case
		r   *http.Request
		u   *user.User
		err bool       // If an error is expected
		exp *user.User // Expected Output Data
	}{
		{
			txt: "nil request",
			err: true,
		},
		{
			txt: "empty request body",
			r:   &http.Request{},
			err: true,
		},
		{
			txt: "empty user",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString("{}")),
			},
			err: true,
		},
		{
			txt: "malformed data",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"id":12}`)),
			},
			u:   &user.User{},
			err: true,
		}, {
			txt: "valid request",
			r: &http.Request{
				Body: ioutil.NopCloser(bytes.NewBuffer(js)),
			},
			u:   &user.User{},
			err: false,
			exp: valid,
		},
	}

	// Loop over different test scenarios
	// tc = Test Case
	for _, tc := range ts {
		// Log the name of the test
		// t.Log function ensures only failing tests display this info.
		t.Log(tc.txt)

		err := bodyToUser(tc.r, tc.u)

		// Test scenario when error is expected
		if tc.err {
			if err == nil {
				t.Error("Expected error, got none")
			}

			// Continue to next test case
			continue
		}

		// Test scenario which should return output and not error
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
			// Continue to next test case
			continue

		}

		// Compare output and expected data
		if !reflect.DeepEqual(tc.u, tc.exp) {
			t.Error("Unmarshalled data is different.")
			t.Error(tc.u)
			t.Error(tc.exp)
		}

	}
}
