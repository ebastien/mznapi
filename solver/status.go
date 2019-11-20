package solver

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
