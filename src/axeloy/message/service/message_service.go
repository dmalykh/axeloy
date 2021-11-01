package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/message/model"
	"github.com/dmalykh/axeloy/axeloy/message/repository"
	"github.com/google/uuid"
)

var (
	ErrUpdateStatus = errors.New(`can't update status`)
	ErrNoMessage    = errors.New(`message does not exist`)
	ErrGetMessage   = errors.New(`message couldn't be got`)
	ErrAddMessage   = errors.New(`message couldn't be added`)
)

func NewMessager(messageRepository repository.MessageRepository) message.Messager {
	return &MessageSevice{
		messageRepository: messageRepository,
	}
}

type MessageSevice struct {
	messageRepository repository.MessageRepository
}

func (m *MessageSevice) Update(ctx context.Context, msg message.Message) error {
	panic("implement me")
}

func (m *MessageSevice) GetById(ctx context.Context, id uuid.UUID) (message.Message, error) {
	msg, err := m.messageRepository.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoMessage) {
			return nil, fmt.Errorf(`%w %s`, ErrNoMessage, id)
		}
		return nil, fmt.Errorf(`%w (%s) %s`, ErrGetMessage, id, err.Error())
	}
	return msg, nil
}

func (m *MessageSevice) Add(ctx context.Context, msg message.Message) error {
	if err := m.messageRepository.Create(ctx, msg); err != nil {
		return fmt.Errorf(`%w %s`, ErrAddMessage, err.Error())
	}
	return nil
}

func (m *MessageSevice) UpdateStatus(ctx context.Context, msg message.Message, status model.Status, info ...string) error {
	if err := m.messageRepository.UpdateStatus(ctx, msg.GetUUID(), status, info...); err != nil {
		return fmt.Errorf(`%w (%s) %s`, ErrUpdateStatus, status, err.Error())
	}
	return nil
}

//var (
//	ErrUnknownState = errors.New(`unknown state`)
//	ErrUpdateState  = errors.New(`can't update state`)
//)
//
//
//status, err := m.convertStateToStatus(state)
//if err != nil {
//	return err
//}
//
//func (m *MessageSevice) convertStateToStatus(state State) (model.Status, error) {
//	switch state {
//	case StateProcessed:
//		return model.StatusProcessed, nil
//	case StateNew:
//		return model.StatusNew, nil
//	case StateSent:
//		return model.StatusSent, nil
//	case StateError:
//		return model.StatusError, nil
//	case StateNoDestinations:
//		return model.StatusNoDestinations, nil
//	}
//	return ``, fmt.Errorf(`%w (%s)`, ErrUnknownState, state)
//}
