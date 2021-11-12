package repository

import (
	"context"
	"errors"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/message/model"
	"github.com/google/uuid"
)

var (
	ErrNoMessage = errors.New(`no message`)
)

type MessageRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*model.Message, error)
	Create(ctx context.Context, msg message.Message) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status, info ...string) error
}
