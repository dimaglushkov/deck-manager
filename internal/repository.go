package internal

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/dimaglushkov/toggl-test-assignment/internal/models"
)

type Repository interface {
	Set(deck models.Deck) error
	Get(uuid uuid.UUID) (*models.Deck, error)
}

type UnknownUUIDError struct {
	id uuid.UUID
}

func NewUnknownUUIDError(id uuid.UUID) error {
	return UnknownUUIDError{id: id}
}

func (e UnknownUUIDError) Error() string {
	return fmt.Sprintf("deck with provided uuid does not exist: %s", e.id)
}
