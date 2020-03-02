package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ebastien/mznapi/solver"
	"github.com/google/uuid"
)

// createHandler loads and compiles a Minizinc model.
func (s *Server) createHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.Lock()
		defer s.lock.Unlock()

		mzn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		model := solver.NewModel(string(mzn))

		fmt.Printf("Compile model: %s\n", model.Minizinc())

		if err := model.Compile(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.models[id] = *model

		w.Header().Set("Location", s.modelURI(id.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

// solveHandler solves the given model and returns the solution as JSON.
func (s *Server) solveHandler(id func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.RLock()
		s.workers <- struct{}{}
		defer func() {
			<-s.workers
			s.lock.RUnlock()
		}()

		uuid, err := uuid.Parse(id(r))

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		model, ok := s.models[uuid]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var solution map[string]interface{}

		fmt.Printf("Solve model: %s\n", model.Flatzinc())

		status, err := model.Solve(&solution)
		if err == nil {
			fmt.Printf("solution = %#v\n", solution)
			fmt.Printf("status = %v\n", status)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m, err := json.Marshal(solution)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = w.Write(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
