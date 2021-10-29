package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	dbmodel "github.com/dmalykh/axeloy/repository/db/router/model"
	"github.com/google/uuid"
	"gopkg.in/reform.v1"
	"strings"
	"time"
)

type TrackRepository struct {
	db *reform.DB
}

var (
	ErrDbQuery        = errors.New(`database error`)
	ErrDbError        = errors.New(`database error`)
	ErrCreateTrack    = errors.New(`track couldn't be saved`)
	ErrUnknownTrack   = errors.New(`track doesn't exists`)
	ErrUpdateTrack    = errors.New(`track wasn't updated`)
	ErrCreateAttempt  = errors.New(`can't create attempt`)
	ErrUnknownAttempt = errors.New(`attempt doesn't exists`)
	ErrUpdateAttempt  = errors.New(`attempt wasn't updated`)
)

func (t *TrackRepository) CreateTrack(ctx context.Context, tracks ...model.Track) error {
	if len(tracks) == 0 {
		return nil
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(`%w %s`, ErrDbError, err.Error())
	}

	var dbtracks = make([]reform.Struct, len(tracks))
	for i, track := range tracks {
		dbtracks[i] = &dbmodel.Track{
			Id:        uuid.New(),
			WayId:     track.GetSenderId(),
			MessageId: track.GetMessageId(),
			Attempts:  0,
			Status:    string(track.GetStatus()), //@TODO: возможно, неправильно из интерфейса статус брать
		}
	}
	if err := tx.InsertMulti(dbtracks...); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%w, rollback error %s`, ErrDbError, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, dbmodel.Source, ErrCreateTrack, err.Error())
	}
	return tx.Commit()
}

func (t *TrackRepository) SetStatus(ctx context.Context, id uuid.UUID, status model.TrackStatus, info ...string) error {
	if err := t.db.UpdateColumns(
		&dbmodel.Track{
			Id:     id,
			Status: string(status),
			Info:   strings.Join(info, "\n"),
		}, `status`); err != nil {
		if errors.Is(err, reform.ErrNoPK) {
			return fmt.Errorf(`%s %w %s`, id.String(), ErrUnknownTrack, err.Error())
		}
		if errors.Is(err, reform.ErrNoRows) {
			return fmt.Errorf(`%s %w %s`, id.String(), ErrUpdateTrack, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, id.String(), ErrDbError, err.Error())
	}
	return nil
}

func (t *TrackRepository) StartAttempt(ctx context.Context, trackId uuid.UUID, status model.AttemptStatus) (*model.Attempt, error) {
	var attempt = &model.Attempt{
		Id:        uuid.New(),
		TrackId:   trackId,
		StartedAt: time.Now(),
		Status:    status,
	}

	if err := t.db.Insert(&dbmodel.Attempt{
		Id:        attempt.Id,
		TrackId:   attempt.TrackId,
		StartedAt: attempt.StartedAt,
		Status:    t.attemptStatusToStr(attempt.Status),
	}); err != nil {
		return nil, fmt.Errorf(`%w %s %+v`, ErrCreateAttempt, err.Error(), attempt)
	}
	return attempt, nil
}

func (t *TrackRepository) FinishAttempt(ctx context.Context, attempt *model.Attempt, status model.AttemptStatus, info ...string) error {
	if err := t.db.UpdateColumns(
		&dbmodel.Attempt{
			Id:     attempt.Id,
			Status: t.attemptStatusToStr(attempt.Status),
			Info:   strings.Join(info, "\n"),
		}, `status`); err != nil {
		if errors.Is(err, reform.ErrNoPK) {
			return fmt.Errorf(`%s %w %s`, attempt.Id.String(), ErrUnknownAttempt, err.Error())
		}
		if errors.Is(err, reform.ErrNoRows) {
			return fmt.Errorf(`%s %w %s`, attempt.Id.String(), ErrUpdateAttempt, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, attempt.Id.String(), ErrDbError, err.Error())
	}
	return nil
}

func (t *TrackRepository) attemptStatusToStr(status model.AttemptStatus) string {
	switch status {
	case model.AttemptStatusInProgress:
		return `in_progress`
	case model.AttemptStatusDone:
		return `done`
	case model.AttemptStatusError:
		return `error`
	}
	return fmt.Sprintf(`unknown (%s)`, status)
}
