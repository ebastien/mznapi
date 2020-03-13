package store

import (
	"fmt"

	"github.com/ebastien/mznapi/solver"
	"github.com/google/uuid"
)

type MemoryStore map[uuid.UUID]*solver.Model

func NewMemoryStore() *MemoryStore {
	s := make(MemoryStore)
	return &s
}

func (s *MemoryStore) Store(id uuid.UUID, model *solver.Model) error {
	(*s)[id] = model
	return nil
}

func (s *MemoryStore) Exists(id uuid.UUID) bool {
	_, ok := (*s)[id]
	return ok
}

func (s *MemoryStore) Load(id uuid.UUID) (*solver.Model, error) {
	m, ok := (*s)[id]
	if !ok {
		return nil, fmt.Errorf("Model not found: %s", id.String())
	}
	return m, nil
}
