package api

import "github.com/ebastien/mznapi/solver"

type ServerState struct {
	model   solver.Model
	workers chan struct{}
}
