package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ebastien/mznapi/solver"
	. "github.com/ebastien/mznapi/testutil"
	"github.com/google/uuid"
)

type ServerError struct {
	msg    string
	status int
}

func (e *ServerError) Error() string { return e.msg }

func serverError(status int, format string, a ...interface{}) error {
	return &ServerError{
		msg:    fmt.Sprintf(format, a...),
		status: status,
	}
}

// newServer creates a server instance for testing.
func newServer() *Server {
	return NewServer("localhost:8080", 1)
}

// createModel posts a new model to a test server instance.
func createModel(handler http.Handler, model string) (string, error) {
	body := strings.NewReader(model)

	req, err := http.NewRequest("POST", "/models", body)
	if err != nil {
		return "", err
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	code := rr.Result().StatusCode
	if code != http.StatusCreated {
		return "", serverError(code, "Expected Created but got %d", code)
	}

	location := ""
	if h := rr.Result().Header["Location"]; len(h) > 0 {
		location = h[0]
	}

	uri, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	return uri.Path, nil
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
		return serverError(code, "Expected OK but got %d", code)
	}

	dec := json.NewDecoder(rr.Body)
	err = dec.Decode(solution)
	if err != nil {
		return err
	}

	return nil
}

func TestCreateHandler(t *testing.T) {
	server := newServer()
	path, err := createModel(server, `var int: variable; constraint variable = 1;`)
	Ok(t, err)
	Assert(t, strings.HasPrefix(path, "/models/"), "Expected redirection but got %s", path)
}

func TestEmptyModel(t *testing.T) {
	server := newServer()
	_, err := createModel(server, "")
	serr := err.(*ServerError)
	Assert(t, serr.status == http.StatusBadRequest, "Expected BadRequest but got %d", serr.status)
}

func TestSolveHandler(t *testing.T) {

	server := newServer()

	uuid := uuid.New()
	model := solver.NewModel("var int: variable; constraint variable = 1;")
	err := model.Compile()
	Ok(t, err)

	server.models[uuid] = *model

	solution := struct{ Variable int }{}

	err = solveModel(server, "/models/"+uuid.String(), &solution)
	Ok(t, err)

	Assert(t, solution.Variable == 1, "Expected solution to be 1 but got %d", solution.Variable)
}

func TestMultipleModels(t *testing.T) {

	server := newServer()

	loc1, err := createModel(server, `var int: variable; constraint variable = 1;`)
	Ok(t, err)

	loc2, err := createModel(server, `var int: variable; constraint variable = 2;`)
	Ok(t, err)

	solution := struct{ Variable int }{}

	err = solveModel(server, loc1, &solution)
	Ok(t, err)
	Assert(t, solution.Variable == 1, "Expected solution to be 1 but got %d", solution.Variable)

	err = solveModel(server, loc2, &solution)
	Ok(t, err)
	Assert(t, solution.Variable == 2, "Expected solution to be 1 but got %d", solution.Variable)
}
