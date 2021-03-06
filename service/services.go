package service

import (
	"github.com/ebastien/mznapi/solver"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// SolverResult contains the output of the solver execution.
type SolverResult struct {
	Status   solver.SolutionStatus
	Solution map[string]interface{}
}

// ModelStore specifies the interface with model stores.
type ModelStore interface {
	Exists(id uuid.UUID) bool
	Store(id uuid.UUID, model *solver.Model) error
	Load(id uuid.UUID) (*solver.Model, error)
}

// CreateModel creates a new model.
func CreateModel(s ModelStore, mzn string) (uuid.UUID, error) {

	model := solver.NewModel(mzn)

	log.Debugf("compiling model: '%.64s'", model.Minizinc())

	if err := model.Compile(); err != nil {
		log.Errorf("Unable to compile the model: %s", err)
		return uuid.Nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Error("Unable to generate a UUID")
		return uuid.Nil, err
	}

	err = s.Store(id, model)
	if err != nil {
		log.Error("Unable to store the model")
		return uuid.Nil, err
	}
	return id, nil
}

// ModelExists checks whether a model already exists.
func ModelExists(s ModelStore, id uuid.UUID) bool {
	return s.Exists(id)
}

// SolveModel executes the solver on an existing model.
func SolveModel(s ModelStore, id uuid.UUID) (*SolverResult, error) {

	model, err := s.Load(id)
	if err != nil {
		log.Errorf("Unable to load the model: %s", err)
		return nil, err
	}

	log.Debugf("solving model: '%.64s'", model.Flatzinc())

	result := SolverResult{}
	result.Status, err = model.Solve(&result.Solution, 50000)
	if err != nil {
		log.Errorf("Unable to solve the model: %s", err)
		return nil, err
	}
	return &result, nil
}
