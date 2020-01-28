package api

import (
	"testing"

	. "github.com/ebastien/mznapi/testutil"
)

func TestPath(t *testing.T) {

	s := NewServer("localhost:8080", 1)
	path := s.Path(ModelResource)

	Assert(t, path == "/models", "Expected path to be /models but got %v", path)
}
