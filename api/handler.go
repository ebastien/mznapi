package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ebastien/mznapi/solver"
)

func initModel(m *solver.Model) {
	m.Init("var int: age; constraint age >= 1; constraint age <= 2; solve satisfy;")

	fmt.Printf("Compile model: %s\n", m.Minizinc())

	if err := m.Compile(); err != nil {
		log.Fatal(err)
	}
}

func (s *ServerState) solveHandler(w http.ResponseWriter, r *http.Request) {
	s.workers <- struct{}{}
	defer func() { <-s.workers }()

	var solution struct{ Age int }

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

func Serve(parallelism int) {
	state := NewState(parallelism)

	initModel(&state.model)

	http.HandleFunc("/", state.solveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
