package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ebastien/hello-go/solver"
)

type HandlerState struct {
	model solver.Model
}

func (s *HandlerState) solveHandler(w http.ResponseWriter, r *http.Request) {

	type Solution struct {
		Age int
	}
	var solution Solution

	fmt.Printf("Solve model: %s\n", s.model.Flatzinc())

	status, err := s.model.Solve(&solution)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("solution = %#v\n", solution)
	fmt.Printf("status = %v\n", status)

	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("Hello, World!")

	var state HandlerState

	state.model.Init("var int: age; constraint age >= 1; constraint age <= 2; solve satisfy;")

	fmt.Printf("Compile model: %s\n", state.model.Minizinc())

	if err := state.model.Compile(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", state.solveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
