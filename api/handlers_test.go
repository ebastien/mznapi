package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ebastien/mznapi/solver"
	. "github.com/ebastien/mznapi/testutil"
	"github.com/google/uuid"
)

type testError struct {
	msg    string
	status int
}

func (e *testError) Error() string { return e.msg }

func newTestError(status int, format string, a ...interface{}) error {
	return &testError{
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

func TestCreateHandler(t *testing.T) {
	server := newServer()
	path, err := createModel(server, `var int: variable; constraint variable = 1;`)
	Ok(t, err)
	Assert(t, strings.HasPrefix(path, "/models/"), "Expected redirection but got %s", path)
}

func TestEmptyModel(t *testing.T) {
	server := newServer()
	_, err := createModel(server, "")
	serr := err.(*testError)
	Assert(t, serr.status == http.StatusBadRequest, "Expected BadRequest but got %d", serr.status)
}

func TestSolveHandler(t *testing.T) {

	server := newServer()

	uuid := uuid.New()
	model := solver.NewModel("var int: variable; constraint variable = 1;")
	err := model.Compile()
	Ok(t, err)

	err = server.models.Store(uuid, model)
	Ok(t, err)

	type SolverSolution struct {
		Variable int
	}
	response := struct {
		Solution SolverSolution         `json:"solution"`
		Status   *solver.SolutionStatus `json:"solver_status"`
	}{}

	err = solveModel(server, "/models/"+uuid.String(), &response)
	Ok(t, err)

	Assert(t, response.Status != nil && *response.Status == solver.SolutionComplete,
		"Expected complete solution")
	Assert(t, response.Solution.Variable == 1,
		"Expected solution to be 1 but got %d", response.Solution.Variable)
}

func TestMultipleModels(t *testing.T) {

	type SolverSolution struct {
		Variable int
	}
	response := struct {
		Solution SolverSolution         `json:"solution"`
		Status   *solver.SolutionStatus `json:"solver_status"`
	}{}

	server := newServer()

	loc1, err := createModel(server, `var int: variable; constraint variable = 1;`)
	Ok(t, err)

	loc2, err := createModel(server, `var int: variable; constraint variable = 2;`)
	Ok(t, err)

	err = solveModel(server, loc1, &response)
	Ok(t, err)
	Assert(t, response.Solution.Variable == 1,
		"Expected solution to be 1 but got %d", response.Solution.Variable)

	err = solveModel(server, loc2, &response)
	Ok(t, err)
	Assert(t, response.Solution.Variable == 2,
		"Expected solution to be 1 but got %d", response.Solution.Variable)
}
