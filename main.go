package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/mitchellh/mapstructure"
)

func compile() string {
	cmd := exec.Command("minizinc",
		"--compile",
		"--input-from-stdin",
		"--output-to-stdout",
		"--no-output-ozn")
	cmd.Stdin = strings.NewReader("var int: age; constraint age >= 1; constraint age <= 2; solve satisfy;")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

// Solver status
const (
	SolutionComplete = iota
	SolutionError
	SolutionUnknown
	SolutionUnbounded
	SolutionUnsatUnbounded
	SolutionUnsat
	SolutionIncomplete
)

func solve(flatzinc string, solution interface{}) int {
	solve := exec.Command("minizinc",
		"--solver", "org.gecode.gecode",
		"--input-from-stdin",
		"--output-mode", "json",
		"--solution-separator", "",
		"--search-complete-msg", fmt.Sprintf(`{ "status": %d }`, SolutionComplete),
		"-a",
	)
	solve.Stdin = strings.NewReader(flatzinc)
	stdout, err := solve.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := solve.Start(); err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(stdout)
	status := SolutionIncomplete
	var doc map[string]interface{}
	for {
		if err := dec.Decode(&doc); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if doc["status"] != nil {
			// Search status document
			status = int(doc["status"].(float64))
		} else {
			// Solution document
			err := mapstructure.Decode(doc, &solution)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if err := solve.Wait(); err != nil {
		log.Fatal(err)
	}
	return status
}

func main() {
	fmt.Println("Hello, World!")

	flatzinc := compile()
	fmt.Print(flatzinc)

	type Solution struct {
		Age int
	}
	var solution Solution

	status := solve(flatzinc, &solution)
	fmt.Printf("solution = %#v\n", solution)
	fmt.Printf("status = %v\n", status)
}
