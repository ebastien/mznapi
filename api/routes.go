package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) bindRoutes() {
	s.router.Route("/models", func(r chi.Router) {
		r.Post("/", s.createHandler(
			// Compose a model resource URI from a model internal identifer.
			func(id string) string {
				return fmt.Sprintf("%s/models/%s", s.baseURL, id)
			},
		))
		r.Get("/{modelID}/solution", s.solveHandler(
			// Retrieve the model internal identifier from the model resource URI.
			func(r *http.Request) string {
				return chi.URLParam(r, "modelID")
			},
		))
	})
}
