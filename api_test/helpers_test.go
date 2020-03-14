package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ebastien/mznapi/api"
	"github.com/ebastien/mznapi/store"
)

// testError contains the details of an error raised by unit tests.
type testError struct {
	msg    string
	status int
}

// testError implement the error interface.
func (e *testError) Error() string { return e.msg }

// newTestError creates a testing error.
func newTestError(status int, format string, a ...interface{}) error {
	return &testError{
		msg:    fmt.Sprintf(format, a...),
		status: status,
	}
}

// newServer creates a server instance for testing.
func newTestServer() *api.Server {
	return api.NewServer("localhost:8080", 1, store.NewMemoryStore())
}

// createModel posts a new model to a test server instance.
func createModel(handler http.Handler, model string) (string, error) {
	body := strings.NewReader(model)

	req, err := http.NewRequest("POST", "/models", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	code := rr.Result().StatusCode
	if code != http.StatusCreated {
		return "", newTestError(code, "Expected Created but got %d", code)
	}

	location, err := rr.Result().Location()
	if err != nil {
		return "", err
	}

	return location.Path, nil
}

// solveModel fetches the solution of a model from a test server instance.
func solveModel(handler http.Handler, path string, solution interface{}) error {
	req, err := http.NewRequest("GET", path+"/solution", nil)
	if err != nil {
		return err
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	code := rr.Result().StatusCode
	if code != http.StatusOK {
		return newTestError(code, "Expected OK but got %d", code)
	}

	contentType := rr.Result().Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "application/json") {
		return newTestError(code, "Expected JSON but got %s", contentType)
	}

	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(solution)
	if err != nil {
		return err
	}

	return nil
}
