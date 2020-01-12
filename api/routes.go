package api

import "net/http"

func matchRoute(method string, path string, handle http.HandlerFunc) (string, http.HandlerFunc) {
	return path, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
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
	s.router.HandleFunc(matchRoute("GET", "/model/solution", s.solveHandler()))
	s.router.HandleFunc(matchRoute("POST", "/model", s.createHandler()))
}
