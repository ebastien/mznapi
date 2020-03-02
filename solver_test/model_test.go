package solver_test

import (
	"testing"

	. "github.com/ebastien/mznapi/testutil"

	"github.com/ebastien/mznapi/solver"
)

func TestCompileOk(t *testing.T) {
	m := solver.NewModel("var int: age;")
	err := m.Compile()
	Ok(t, err)
	fzn := m.Flatzinc()
	Assert(t, len(fzn) != 0, "Compilation expected to generate non-empty flatzinc model")
}

func TestCompileFail(t *testing.T) {
	m := solver.NewModel("invalid model")
	err := m.Compile()
	Assert(t, err != nil, "Compilation expected to fail but got %v", err)
}

func TestSolveComplete(t *testing.T) {
	m := solver.NewModel(`var int: age; constraint age >= 2 /\ age <= 4; solve maximize age;`)
	err := m.Compile()
	Ok(t, err)
	solution := struct{ Age int }{}
	status, err := m.Solve(&solution)
	Ok(t, err)
	Assert(t, status == solver.SolutionComplete, "Solution expected to be complete but got %v", status)
	Assert(t, solution.Age == 4, "Solution expected to be 4 but got %v", solution.Age)
}

func TestSolveUnsat(t *testing.T) {
	m := solver.NewModel(`var int: age; constraint age < 2; constraint age > 4; solve satisfy;`)
	err := m.Compile()
	Ok(t, err)
	solution := struct{}{}
	status, err := m.Solve(&solution)
	Ok(t, err)
	Assert(t, status == solver.SolutionUnsat, "Solution expected to be unsatisfiable but got %v", status)
}
