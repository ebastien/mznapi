package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// accept validates the content type of a request
func accept(contentType string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), contentType) {
			w.Header().Set("Accept", contentType)
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		h(w, r)
	}
}

// newRouter creates an http.Handler routing resources to the appropriate handler.
func newRouter(s *Server) http.Handler {
	router := chi.NewRouter()
	router.Route("/models", func(r chi.Router) {
		r.Post("/", accept("application/json", s.createHandler(
			// Compose a model resource URI from a model internal identifer.
			func(id string) string {
				return fmt.Sprintf("%s/models/%s", s.baseURL, id)
			},
		)))
		r.Get("/{modelID}/solution", s.solveHandler(
			// Retrieve the model internal identifier from the model resource URI.
			func(r *http.Request) string {
				return chi.URLParam(r, "modelID")
			},
		))
	})
	return router
}
