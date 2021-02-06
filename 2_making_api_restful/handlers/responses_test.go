package handlers

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/nitishgalaxy/go-rest-api/user"
	"gopkg.in/mgo.v2/bson"
)

type response struct {
	header http.Header
	code   int
	body   []byte
}

const (
	dbpath = "users.db"
)

// Copied these from writer_test.go as we cannot share tests between packages.
type mockWriter response

func newMockWriter() *mockWriter {
	return &mockWriter{
		body:   []byte{},
		header: http.Header{},
	}
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.body = make([]byte, len(b))
	for k, v := range b {
		mw.body[k] = v
	}
	return len(b), nil
}

func (mw *mockWriter) WriteHeader(code int) {
	mw.code = code
}

func (mw *mockWriter) Header() http.Header {
	return mw.header
}

func TestMain(m *testing.M) {
	m.Run() // Execute all tests for given package
	os.Remove(dbpath)
}

func prepDB(n int) error {
	os.Remove(dbpath)

	for i := 0; i < n; i++ {
		u := &user.User{
			ID:   bson.NewObjectId(),
			Name: "John_" + strconv.Itoa(i),
			Role: "Tester",
		}

		err := u.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func makeRequest() (*http.Request, error) {
	res := "/users"
	u, err := url.Parse(res)

	if err != nil {
		return nil, err
	}

	return &http.Request{
		URL:    u,
		Header: http.Header{},
		Method: http.MethodGet,
	}, nil

}

func getAll(b *testing.B, r *http.Request) {
	prepDB(100)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mw := newMockWriter()
		b.StartTimer()
		UsersRouter(mw, r)
	}
}

func BenchmarkGetAllNonCached(b *testing.B) {
	r, err := makeRequest()
	if err != nil {
		b.Fatal(err)
	}
	r.Header.Add("Cache-Control", "no-cache")
	getAll(b, r)
}

func BenchmarkGetCached(b *testing.B) {
	r, err := makeRequest()
	if err != nil {
		b.Fatal(err)
	}

	getAll(b, r)
}
