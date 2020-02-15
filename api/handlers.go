package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// createHandler loads and compiles a Minizinc model.
func (s *Server) createHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.Lock()
		defer s.lock.Unlock()

		mzn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			s.model.Init(string(mzn))

			fmt.Printf("Compile model: %s\n", s.model.Minizinc())

			if err := s.model.Compile(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Set("Location", s.modelURI("1"))
				w.WriteHeader(http.StatusCreated)
			}
		}
	}
}

// solveHandler solves the current model and returns the solution as JSON.
func (s *Server) solveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.RLock()
		s.workers <- struct{}{}
		defer func() {
			<-s.workers
			s.lock.RUnlock()
		}()

		var solution map[string]interface{}

		fmt.Printf("Solve model: %s\n", s.model.Flatzinc())

		status, err := s.model.Solve(&solution)
		if err == nil {
			fmt.Printf("solution = %#v\n", solution)
			fmt.Printf("status = %v\n", status)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			m, err := json.Marshal(solution)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				_, err := w.Write(m)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.WriteHeader(http.StatusOK)
			}
		}
	}
}
