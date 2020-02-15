package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/ebastien/mznapi/solver"
	"github.com/go-chi/chi"
)

// Server maintains the state of the HTTP APIs.
type Server struct {
	model   solver.Model
	address string
	baseURL string
	router  *chi.Mux
	lock    sync.RWMutex
	workers chan struct{}
}

// NewServer creates a new server instance.
func NewServer(addr string, parallelism int) *Server {
	return &Server{
		address: addr,
		baseURL: "http://" + addr,
		router:  chi.NewRouter(),
		workers: make(chan struct{}, parallelism),
	}
}

// Serve runs the server main loop to handle incoming connections.
func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(s.address, s.router))
}

// ServeHTTP proxies to the underlying router implementation.
// It implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
