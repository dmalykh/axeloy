package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/dmalykh/axeloy/axeloy/router/repository"
	"github.com/dmalykh/axeloy/axeloy/way"
)

var (
	ErrUpdateTrack = errors.New(`can't update track`)
	ErrGetSender   = errors.New(`can't get sender`)
	ErrGetMessage  = errors.New(`can't get message`)
	ErrSendMessage = errors.New(`can't sent message`)
	ErrSaveAttempt = errors.New(`can't save attempt`)
)

type TrackService struct {
	trackRepository repository.TrackRepository
	wayService      way.Wayer
	messageService  message.Messager
}

func (t *TrackService) DefineTracks(ctx context.Context, m message.Message, destination Destination) ([]*model.Track, error) {
	var plannedTracks = make([]*model.Track, 0)
	var errorTracks = make([]*model.Track, 0)
	for _, w := range destination.GetWays(ctx) {
		var track = &model.Track{
			SenderId:  w.GetId(),
			MessageId: m.GetUUID(),
			Attempts:  0,
			Profile:   destination.GetProfile(ctx),
			Status:    model.Planned,
		}
		errinfo, err := w.ValidateProfile(ctx, destination.GetProfile(ctx))
		if err != nil {
			track.Status = model.Error
			track.Info = func(info map[string]string) string {
				data, _ := json.Marshal(info)
				return string(data)
			}(errinfo)
			errorTracks = append(plannedTracks, track)
		} else {
			plannedTracks = append(plannedTracks, track)
		}
	}
	if err := t.trackRepository.CreateTrack(ctx, append(plannedTracks, errorTracks...)...); err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrCreateTrack, err.Error())
	}
	return plannedTracks, nil
}

func (t *TrackService) GetTracks(ctx context.Context, m message.Message) ([]*model.Track, error) {
	tracks, err := t.trackRepository.GetTracksByMessageId(ctx, m.GetUUID())
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrGetTrack, err.Error())
	}
	return tracks, nil
}

func (t *TrackService) Send(ctx context.Context, track Track) error {
	//Start attempt to send message
	attempt, err := t.trackRepository.StartAttempt(ctx, track.GetId())
	if err != nil {
		return fmt.Errorf(`%w %s`, ErrUpdateTrack, err.Error())
	}

	//Set track's status in progress...
	if err := t.trackRepository.SetStatus(ctx, track.GetId(), model.Process); err != nil {
		if err := t.trackRepository.FinishAttempt(ctx, attempt, model.AttemptStatusError, ErrUpdateTrack.Error(), err.Error()); err != nil {
			return fmt.Errorf(`%w %s`, ErrSaveAttempt, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrUpdateTrack, err.Error())
	}

	//Get sender
	sender, err := t.wayService.GetSenderById(ctx, track.GetSenderId())
	if err != nil {
		if err := t.trackRepository.FinishAttempt(ctx, attempt, model.AttemptStatusError, ErrGetSender.Error(), err.Error()); err != nil {
			return fmt.Errorf(`%w %s`, ErrSaveAttempt, err.Error())
		}
		return fmt.Errorf(`%w %s %s`, ErrGetSender, track.GetSenderId().String(), err.Error())
	}

	//Get message
	m, err := t.messageService.GetById(ctx, track.GetMessageId())
	if err != nil {
		if err := t.trackRepository.FinishAttempt(ctx, attempt, model.AttemptStatusError, ErrGetMessage.Error(), err.Error()); err != nil {
			return fmt.Errorf(`%w %s`, ErrSaveAttempt, err.Error())
		}
		return fmt.Errorf(`%w %s %s`, ErrGetMessage, track.GetMessageId().String(), err.Error())
	}

	//Send it!
	info, err := sender.Send(ctx, track.GetProfile(), m)
	if err != nil {
		if err := t.trackRepository.FinishAttempt(ctx, attempt, model.AttemptStatusError, append(info, ErrSaveAttempt.Error(), err.Error())...); err != nil {
			return fmt.Errorf(`%w %s`, ErrSaveAttempt, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrSendMessage, err.Error())
	}

	//All done!
	if err := t.trackRepository.FinishAttempt(ctx, attempt, model.AttemptStatusDone, info...); err != nil {
		return fmt.Errorf(`%w %s`, ErrSaveAttempt, err.Error())
	}
	if err := t.trackRepository.SetStatus(ctx, track.GetId(), model.Done); err != nil {
		return fmt.Errorf(`%w %s`, ErrUpdateTrack, err.Error())
	}

	return nil
}
