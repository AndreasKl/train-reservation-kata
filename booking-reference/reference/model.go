package reference

import "github.com/google/uuid"

type ID uuid.UUID

type Reference struct {
	ID ID
}
