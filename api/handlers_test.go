package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ebastien/mznapi/solver"
	. "github.com/ebastien/mznapi/testutil"
	"github.com/google/uuid"
)

func TestCreateHandler(t *testing.T) {

	server := NewServer("localhost:8080", 1)
	server.routes()

	body := strings.NewReader(`var int: age; constraint age = 1;`)

	req, err := http.NewRequest("POST", "/models", body)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	code := rr.Result().StatusCode
	redirect := ""
	if h := rr.Result().Header["Location"]; len(h) > 0 {
		redirect = h[0]
	}

	Assert(t, code == http.StatusCreated, "Expected Created but got %v", code)
	Assert(t,
		strings.HasPrefix(redirect, "http://localhost:8080/models/"),
		"Expected redirection but got %v", redirect)
}

func TestSolveHandler(t *testing.T) {

	server := NewServer("localhost:8080", 1)
	server.routes()

	uuid := uuid.New()
	model := solver.NewModel("var int: variable; constraint variable = 1;")
	err := model.Compile()
	Ok(t, err)

	server.models[uuid] = *model

	req, err := http.NewRequest("GET", "/models/"+uuid.String()+"/solution", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusOK, "Expected OK but got %v", rr.Code)

	dec := json.NewDecoder(rr.Body)
	solution := struct{ Variable int }{}
	err = dec.Decode(&solution)
	Ok(t, err)

	Assert(t, solution.Variable == 1, "Expected solution to be 1 but got %v", solution.Variable)
}
