package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/ebastien/mznapi/solver"
)

type Server struct {
	model   solver.Model
	address string
	baseURL string
	router  *http.ServeMux
	lock    sync.RWMutex
	workers chan struct{}
}

func NewServer(addr string, parallelism int) *Server {
	return &Server{
		address: addr,
		baseURL: "http://" + addr,
		router:  http.NewServeMux(),
		workers: make(chan struct{}, parallelism),
	}
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(s.address, s.router))
}
