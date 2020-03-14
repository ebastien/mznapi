package api_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ebastien/mznapi/api"
	"github.com/ebastien/mznapi/solver"
	"github.com/ebastien/mznapi/store"
	. "github.com/ebastien/mznapi/testutil"
	"github.com/google/uuid"
)

func TestCreateHandler(t *testing.T) {
	server := newTestServer()
	path, err := createModel(server, `var int: variable; constraint variable = 1;`)
	Ok(t, err)
	Assert(t, strings.HasPrefix(path, "/models/"), "Expected redirection but got %s", path)
}

func TestEmptyModel(t *testing.T) {
	server := newTestServer()
	_, err := createModel(server, "")
	serr := err.(*testError)
	Assert(t, serr.status == http.StatusBadRequest, "Expected BadRequest but got %d", serr.status)
}

func TestSolveHandler(t *testing.T) {

	store := store.NewMemoryStore()
	server := api.NewServer("localhost:8080", 1, store)

	uuid := uuid.New()
	model := solver.NewModel("var int: variable; constraint variable = 1;")
	err := model.Compile()
	Ok(t, err)

	err = store.Store(uuid, model)
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

	server := newTestServer()

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
