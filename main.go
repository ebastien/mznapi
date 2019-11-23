package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebastien/hello-go/solver"
)

type ServerState struct {
	model   solver.Model
	workers chan struct{}
}

func initModel(m *solver.Model) {
	m.Init("var int: age; constraint age >= 1; constraint age <= 2; solve satisfy;")

	fmt.Printf("Compile model: %s\n", m.Minizinc())

	if err := m.Compile(); err != nil {
		log.Fatal(err)
	}
}

func solveModel(m solver.Model) error {
	type Solution struct {
		Age int
	}
	var solution Solution

	fmt.Printf("Solve model: %s\n", m.Flatzinc())

	status, err := m.Solve(&solution)
	if err == nil {
		fmt.Printf("solution = %#v\n", solution)
		fmt.Printf("status = %v\n", status)
	}
	return err
}

func (s *ServerState) solveHandler(w http.ResponseWriter, r *http.Request) {
	s.workers <- struct{}{}
	defer func() { <-s.workers }()

	// time.Sleep(5 * time.Second)

	err := solveModel(s.model)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	fmt.Println("Hello, World!")

	state := ServerState{
		workers: make(chan struct{}, 3),
	}

	initModel(&state.model)

	http.HandleFunc("/", state.solveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
