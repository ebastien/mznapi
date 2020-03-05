package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ebastien/mznapi/solver"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// createHandler loads and compiles a Minizinc model.
func (s *Server) createHandler(uri func(id string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.Lock()
		defer s.lock.Unlock()

		mzn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Unable to read the request body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(mzn) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		model := solver.NewModel(string(mzn))

		log.Debugf("compiling model: '%.64s'", model.Minizinc())

		if err := model.Compile(); err != nil {
			log.Errorf("Unable to compile the model: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			log.Error("Unable to generate a UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.models[id] = *model

		w.Header().Set("Location", uri(id.String()))
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

		log.Debugf("solving model: '%.64s'", model.Flatzinc())

		status, err := model.Solve(&solution, 50000)
		if err == nil {
			log.WithField("status", status).Debugf("solution: %#v", solution)
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
