package solver_test

import (
	"testing"

	. "github.com/ebastien/mznapi/testutil"

	"github.com/ebastien/mznapi/solver"
)

func TestCompileOk(t *testing.T) {
	m := solver.Model{}
	m.Init("var int: age;")
	err := m.Compile()
	Ok(t, err)
}

func TestCompileFail(t *testing.T) {
	m := solver.Model{}
	m.Init("invalid model")
	err := m.Compile()
	Assert(t, err != nil, "Compilation expected to fail but got %v", err)
}
