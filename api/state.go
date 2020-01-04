package api

import (
	"net/http"

	"github.com/ebastien/mznapi/solver"
)

type serverState struct {
	model   solver.Model
	router  *http.ServeMux
	workers chan struct{}
}

func newState(parallelism int) *serverState {
	return &serverState{
		router:  http.NewServeMux(),
		workers: make(chan struct{}, parallelism),
	}
}
