package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/ebastien/mznapi/testutil"
)

func TestSolveHandler(t *testing.T) {

	state := NewState(1)
	state.model.Init("solve satisfy;")
	err := state.model.Compile()
	Ok(t, err)

	req, err := http.NewRequest("GET", "/", nil)
	Ok(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(state.solveHandler)

	handler.ServeHTTP(rr, req)

	Assert(t, rr.Code == http.StatusOK, "Expected OK but got %v", rr.Code)
}
