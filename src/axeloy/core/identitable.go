package core

import "github.com/google/uuid"

type Identitable interface {
	GetId() uuid.UUID
}
