package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/ebastien/mznapi/testutil"
)

func TestCreateHandler(t *testing.T) {

	server := NewServer("localhost:8080", 1)

	body := strings.NewReader(`var int: age; constraint age = 1;`)

	req, err := http.NewRequest("POST", "/models", body)
	Ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.createHandler())

	handler.ServeHTTP(rr, req)

	code := rr.Result().StatusCode
	redirect := ""
	if h := rr.Result().Header["Location"]; len(h) > 0 {
		redirect = h[0]
	}

	Assert(t, code == http.StatusCreated, "Expected Created but got %v", code)
	Assert(t, redirect == "http://localhost:8080/models/1", "Expected redirection but got %v", redirect)
}

func TestSolveHandler(t *testing.T) {

	server := NewServer("localhost:8080", 1)
	server.model.Init("var int: age; constraint age = 1;")
	err := server.model.Compile()
	Ok(t, err)

	req, err := http.NewRequest("GET", "/", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.solveHandler())

	handler.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusOK, "Expected OK but got %v", rr.Code)

	dec := json.NewDecoder(rr.Body)
	solution := struct{ Age int }{}
	err = dec.Decode(&solution)
	Ok(t, err)

	Assert(t, solution.Age == 1, "Expected solution to be 1 but got %v", solution.Age)
}
