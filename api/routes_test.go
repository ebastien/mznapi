package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ebastien/mznapi/testutil"
)

func TestUnknownURL(t *testing.T) {
	server := NewServer(":8080", 1)
	server.routes()

	req, err := http.NewRequest("GET", "/model/unknown_url", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusNotFound, "Expected NotFound but got %v", rr.Code)
}

func TestInvalidMethod(t *testing.T) {
	server := NewServer(":8080", 1)
	server.routes()

	req, err := http.NewRequest("POST", "/model/solution", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusMethodNotAllowed, "Expected MethodNotAllowed but got %v", rr.Code)
}
