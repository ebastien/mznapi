package store

import (
	"fmt"

	"github.com/ebastien/mznapi/solver"
	"github.com/google/uuid"
)

// MemoryStore maintains an in-memory store of models.
type MemoryStore map[uuid.UUID]*solver.Model

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	s := make(MemoryStore)
	return &s
}

// Store registers a model in the in-memory store.
// It implements the service.ModelStore interface.
func (s *MemoryStore) Store(id uuid.UUID, model *solver.Model) error {
	(*s)[id] = model
	return nil
}

// Exists checks whether a model exists in the in-memory store.
// It implements the service.ModelStore interface.
func (s *MemoryStore) Exists(id uuid.UUID) bool {
	_, ok := (*s)[id]
	return ok
}

// Load fetches a model from the in-memory store.
// It implements the service.ModelStore interface.
func (s *MemoryStore) Load(id uuid.UUID) (*solver.Model, error) {
	m, ok := (*s)[id]
	if !ok {
		return nil, fmt.Errorf("Model not found: %s", id.String())
	}
	return m, nil
}
