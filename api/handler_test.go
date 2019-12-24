package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ebastien/mznapi/testutil"
)

func TestSolveHandler(t *testing.T) {

	state := NewState(1)
	state.model.Init("var int: age; constraint age = 1;")
	err := state.model.Compile()
	Ok(t, err)

	req, err := http.NewRequest("GET", "/", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(state.solveHandler)

	handler.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusOK, "Expected OK but got %v", rr.Code)

	dec := json.NewDecoder(rr.Body)
	solution := struct{ Age int }{}
	err = dec.Decode(&solution)
	Ok(t, err)

	Assert(t, solution.Age == 1, "Expected solution to be 1 but got %v", solution.Age)
}
