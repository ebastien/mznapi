package api

import (
	"fmt"

	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.router.Route("/models", func(r chi.Router) {
		r.Post("/", s.createHandler())
		r.Get("/{modelID}/solution", s.solveHandler())
	})
}

func (s *Server) modelURI(id string) string {
	return fmt.Sprintf("%s/models/%s", s.baseURL, id)
}
