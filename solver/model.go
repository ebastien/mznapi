package solver

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Model represents a Minizinc model.
type Model struct {
	mzn string
	fzn string
}

// NewModel returns a pointer to a new model.
func NewModel(mzn string) *Model {
	return &Model{
		mzn: mzn,
	}
}

// Minizinc returns the original Minizinc representation of the model.
func (m *Model) Minizinc() string {
	return m.mzn
}

// Flatzinc returns the compiled Flatzinc representation of the model.
func (m *Model) Flatzinc() string {
	return m.fzn
}

// Compile translates a model to its Flatzinc representation.
func (m *Model) Compile() error {
	if len(m.mzn) == 0 {
		return fmt.Errorf("Model not initialized")
	}
	cmd := exec.Command("minizinc",
		"--compile",
		"--input-from-stdin",
		"--output-to-stdout",
		"--no-output-ozn")
	cmd.Stdin = strings.NewReader(m.mzn)
	out, err := cmd.Output()
	if err == nil {
		m.fzn = string(out)
	}
	return err
}

// Solve runs the solver on the compiled model and tries to retrieve a solution.
func (m *Model) Solve(solution interface{}, timeout int) (SolutionStatus, error) {
	status := SolutionIncomplete
	solve := exec.Command("minizinc",
		"--solver", "org.gecode.gecode",
		"--input-from-stdin",
		"--output-mode", "json",
		"--solution-separator", "",
		"--search-complete-msg", fmt.Sprintf(`{ "status": %d }`, SolutionComplete),
		"--unsatisfiable-msg", fmt.Sprintf(`{ "status": %d }`, SolutionUnsat),
		"--unbounded-msg", fmt.Sprintf(`{ "status": %d }`, SolutionUnbounded),
		"--unsatorunbnd-msg", fmt.Sprintf(`{ "status": %d }`, SolutionUnsatUnbounded),
		"--unknown-msg", fmt.Sprintf(`{ "status": %d }`, SolutionUnknown),
		"--error-msg", fmt.Sprintf(`{ "status": %d }`, SolutionError),
		"--time-limit", strconv.Itoa(timeout),
		"-a",
	)
	solve.Stdin = strings.NewReader(m.fzn)
	stdout, err := solve.StdoutPipe()
	if err != nil {
		return status, err
	}
	if err := solve.Start(); err != nil {
		return status, err
	}
	dec := json.NewDecoder(stdout)
	var doc map[string]interface{}
	for {
		if err := dec.Decode(&doc); err == io.EOF {
			break
		} else if err != nil {
			return status, err
		}
		if doc["status"] != nil {
			// Search status JSON document
			status = SolutionStatus(doc["status"].(float64))
		} else {
			// Solution JSON document
			err := mapstructure.Decode(doc, solution)
			if err != nil {
				return status, err
			}
		}
	}
	err = solve.Wait()
	return status, err
}
