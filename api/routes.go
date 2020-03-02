package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Route("/models", func(r chi.Router) {
		r.Post("/", s.createHandler())
		r.Get("/{modelID}/solution", s.solveHandler(
			func(r *http.Request) string { return chi.URLParam(r, "modelID") },
		))
	})
}

func (s *Server) modelURI(id string) string {
	return fmt.Sprintf("%s/models/%s", s.baseURL, id)
}
