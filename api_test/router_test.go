package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/ebastien/mznapi/testutil"
)

func TestUnknownURL(t *testing.T) {
	server := newTestServer()

	req, err := http.NewRequest("GET", "/models/unknown_url", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusNotFound, "Expected NotFound but got %v", rr.Code)
}

func TestInvalidMethod(t *testing.T) {
	server := newTestServer()

	req, err := http.NewRequest("POST", "/models/1/solution", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusMethodNotAllowed,
		"Expected MethodNotAllowed but got %v", rr.Code)
}

func TestInvalidContentType(t *testing.T) {
	server := newTestServer()

	body := strings.NewReader("some invalid content")
	req, err := http.NewRequest("POST", "/models", body)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusUnsupportedMediaType,
		"Expected UnsupportedMediaType but got %v", rr.Code)

	accept := rr.Result().Header.Get("Accept")

	Assert(t, strings.HasPrefix(accept, "application/json"),
		"Expected application/json but got %s", accept)
}
