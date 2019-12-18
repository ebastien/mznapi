package solver

// SolutionStatus qualifies the solution returned by the solver.
type SolutionStatus int

const (
	SolutionComplete SolutionStatus = iota
	SolutionError
	SolutionUnknown
	SolutionUnbounded
	SolutionUnsatUnbounded
	SolutionUnsat
	SolutionIncomplete
)
