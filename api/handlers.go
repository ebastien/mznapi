package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ebastien/mznapi/service"
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
			log.Errorf("Unable to read the request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(mzn) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		uuid, err := service.CreateModel(s.models, string(mzn))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", uri(uuid.String()))
		w.WriteHeader(http.StatusCreated)
	}
}

// solveHandler solves the given model and returns the solution as JSON.
func (s *Server) solveHandler(id func(*http.Request) string) http.HandlerFunc {

	type Response struct {
		Solution map[string]interface{} `json:"solution"`
		Status   solver.SolutionStatus  `json:"solver_status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		s.lock.RLock()
		s.workers <- struct{}{}
		defer func() {
			<-s.workers
			s.lock.RUnlock()
		}()

		uuid, err := uuid.Parse(id(r))
		if err != nil {
			log.Errorf("Unable to parse the model uuid: %s", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if !service.ModelExists(s.models, uuid) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		res, err := service.SolveModel(s.models, uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := Response{
			Solution: res.Solution,
			Status:   res.Status,
		}
		msg, err := json.Marshal(response)
		if err != nil {
			log.Errorf("Unable to serialize the response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(msg)
		if err != nil {
			log.Errorf("Unable to write the response: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
