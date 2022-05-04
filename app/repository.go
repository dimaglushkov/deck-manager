package app

import (
	"fmt"
	"github.com/dimaglushkov/toggl-test-assignment/app/models"
	"github.com/google/uuid"
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
