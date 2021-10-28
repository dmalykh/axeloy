package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	dbmodel "github.com/dmalykh/axeloy/repository/db/router/model"
	"github.com/google/uuid"
	"gopkg.in/reform.v1"
)

type TrackRepository struct {
	db *reform.DB
}

var ErrDbQuery = errors.New(`database error`)

var ErrDbError = errors.New(`database error`)
var ErrCreateTrack = errors.New(`track couldn't be saved`)

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

func (t *TrackRepository) GetTracksByMessageId(ctx context.Context, messageId uuid.UUID) ([]*model.Track, error) {
	panic("implement me")
}
