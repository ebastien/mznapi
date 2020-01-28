package api

import (
	"net/http"
)

type PathMatcher interface {
	Route() string
	Match(path string) bool
}

type ExactMatch string

func (pattern ExactMatch) Match(path string) bool {
	return string(pattern) == path
}

func (pattern ExactMatch) Route() string {
	return string(pattern)
}

func matchRoute(method string, pattern PathMatcher, handle http.HandlerFunc) (string, http.HandlerFunc) {
	return pattern.Route(), func(w http.ResponseWriter, r *http.Request) {
		if !pattern.Match(r.URL.Path) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handle(w, r)
	}
}

func (s *Server) routes() {
	s.router.HandleFunc(matchRoute("GET", ExactMatch(s.Path(ModelResource, "1", "solution")), s.solveHandler()))
	s.router.HandleFunc(matchRoute("POST", ExactMatch(s.Path(ModelResource)), s.createHandler()))
}
