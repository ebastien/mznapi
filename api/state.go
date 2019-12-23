package api

import "github.com/ebastien/mznapi/solver"

type ServerState struct {
	model   solver.Model
	workers chan struct{}
}

func NewState(parallelism int) *ServerState {
	return &ServerState{
		workers: make(chan struct{}, parallelism),
	}
}
