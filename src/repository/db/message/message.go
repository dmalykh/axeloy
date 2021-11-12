package message

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/message/location"
	"github.com/dmalykh/axeloy/axeloy/message/model"
	"github.com/dmalykh/axeloy/axeloy/message/payload"
	"github.com/dmalykh/axeloy/axeloy/profile"
	profilemodel "github.com/dmalykh/axeloy/axeloy/profile/model"
	dbmodel "github.com/dmalykh/axeloy/repository/db/message/model"
	"github.com/google/uuid"
	"gopkg.in/reform.v1"
	"strings"
	"time"
)

var (
	ErrNoMessage              = errors.New(`no message`)
	ErrNoPayload              = errors.New(`no messages payload`)
	ErrNoProfile              = errors.New(`no messages profile`)
	ErrUnknownStructReturned  = errors.New(`returned unknown struct`)
	ErrUnknownStatus          = errors.New(`returned unknown status`)
	ErrCreateMessage          = errors.New(`can't save message`)
	ErrSavePayload            = errors.New(`can't save payload`)
	ErrSaveSourceProfile      = errors.New(`can't save source profile`)
	ErrSaveDestinationProfile = errors.New(`can't save destination profile`)
	ErrDbUpdate               = errors.New(`can't update, db error`)
)

type MessageRepository struct {
	db *reform.DB
}

func NewMessageRepository(db *reform.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (m *MessageRepository) GetById(ctx context.Context, id uuid.UUID) (*model.Message, error) {
	var message model.Message

	//Get message
	{
		row, err := m.db.FindOneFrom(dbmodel.MessageTable, "id", id)
		if err != nil {
			return nil, fmt.Errorf(`%w %s %s`, ErrNoMessage, id.String(), err.Error())
		}
		msg, ok := row.(*dbmodel.Message)
		if !ok {
			return nil, fmt.Errorf(`%w %+v`, ErrUnknownStructReturned, row)
		}
		message.Id = msg.Id
		message.Info = strings.Split(msg.Info, "\t")
		message.Source = &model.Location{Way: msg.SourceWayName}
		message.Status, err = m.strToStatus(msg.Status)
		if err != nil {
			message.Info = []string{err.Error()}
		}
	}

	//Get payload
	{
		message.Payload = make(payload.Payload)
		rows, err := m.db.FindAllFrom(dbmodel.PayloadTable, "message_id", message.Id)
		if err != nil {
			return nil, fmt.Errorf(`%w %s %s`, ErrNoPayload, message.Id.String(), err.Error())
		}
		for _, row := range rows {
			p, ok := row.(*dbmodel.Payload)
			if !ok {
				return nil, fmt.Errorf(`%w %+v`, ErrUnknownStructReturned, row)
			}
			if _, exists := message.Payload[p.Key]; !exists {
				message.Payload[p.Key] = make([]string, 0)
			}
			message.Payload[p.Key] = append(message.Payload[p.Key], p.Value)
		}
	}

	//Get profiles
	{
		var sourceProfile = make(profile.Fields)
		var destinations = make(map[string]profile.Fields)
		rows, err := m.db.FindAllFrom(dbmodel.ProfileTable, "message_id", message.Id)
		if err != nil {
			return nil, fmt.Errorf(`%w %s %s`, ErrNoProfile, message.Id.String(), err.Error())
		}
		for _, row := range rows {
			p, ok := row.(*dbmodel.Profile)
			if !ok {
				return nil, fmt.Errorf(`%w %+v`, ErrUnknownStructReturned, row)
			}
			//Source profile hasn't destination's way name
			if p.DestinationWayName == nil {
				if _, exists := sourceProfile[p.Key]; !exists {
					sourceProfile[p.Key] = make([]string, 0)
				}
				sourceProfile[p.Key] = append(sourceProfile[p.Key], p.Value)
			} else {
				if _, exists := destinations[*p.DestinationWayName][p.Key]; !exists {
					destinations[*p.DestinationWayName][p.Key] = make([]string, 0)
				}
				destinations[*p.DestinationWayName][p.Key] = append(destinations[*p.DestinationWayName][p.Key], p.Value)
			}
		}
		//Set profile data
		message.Source = &model.Location{
			Way: message.GetSource().GetWay(),
			Profile: &profilemodel.Profile{
				Fields: sourceProfile,
			},
		}
		//Set destinations
		message.Destination = make([]location.Location, 0)
		for wayName, fields := range destinations {
			message.Destination = append(message.Destination, &model.Location{
				Way: wayName,
				Profile: &profilemodel.Profile{
					Fields: fields,
				},
			})
		}
	}

	return &message, nil
}

func (m *MessageRepository) Create(ctx context.Context, msg message.Message) error {
	msg.SetId(uuid.New())

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err //@TODO
	}

	//Save message
	{
		if err := tx.Insert(&dbmodel.Message{
			Id:            msg.GetId(),
			CreatedAt:     time.Now(),
			SourceWayName: msg.GetSource().GetWay(),
			Info:          strings.Join(msg.GetInfo(), "\t"),
			Status:        m.statusToSrt(msg.GetStatus()),
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf(`%w %s, rollback error %s`, ErrCreateMessage, msg.GetId().String(), err.Error())
			}
			return fmt.Errorf(`%w %s`, ErrCreateMessage, err.Error())
		}
	}

	//Add payload
	{
		var payload = make([]reform.Struct, 0)
		for key, values := range msg.GetPayload() {
			for _, value := range values {
				payload = append(payload, &dbmodel.Payload{
					Id:        uuid.New(),
					MessageId: msg.GetId(),
					Key:       key,
					Value:     value,
				})
			}
		}
		if err := tx.InsertMulti(payload...); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf(`%w %s, rollback error %s`, ErrSavePayload, msg.GetId().String(), err.Error())
			}
			return fmt.Errorf(`%w %s`, ErrSavePayload, err.Error())
		}
	}

	//Save source profile
	{
		var profiles = make([]reform.Struct, 0)
		for key, values := range msg.GetSource().GetProfile().GetFields() {
			for _, value := range values {
				profiles = append(profiles, &dbmodel.Profile{
					Id:        uuid.New(),
					MessageId: msg.GetId(),
					Key:       key,
					Value:     value,
				})
			}
		}
		if err := tx.InsertMulti(profiles...); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf(`%w %s, rollback error %s`, ErrSaveSourceProfile, msg.GetId().String(), err.Error())
			}
			return fmt.Errorf(`%w %s`, ErrSaveSourceProfile, err.Error())
		}
	}

	//Save destinations
	{
		var profiles = make([]reform.Struct, 0)
		for _, destination := range msg.GetDestinations() {
			var wayName = destination.GetWay()
			for key, values := range destination.GetProfile().GetFields() {
				for _, value := range values {
					profiles = append(profiles, &dbmodel.Profile{
						Id:                 uuid.New(),
						MessageId:          msg.GetId(),
						DestinationWayName: &wayName,
						Key:                key,
						Value:              value,
					})
				}
			}
		}
		if err := tx.InsertMulti(profiles...); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf(`%w %s, rollback error %s`, ErrSaveDestinationProfile, msg.GetId().String(), err.Error())
			}
			return fmt.Errorf(`%w %s`, ErrSaveDestinationProfile, err.Error())
		}
	}

	return tx.Commit()
}

func (m *MessageRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status, info ...string) error {
	if err := m.db.UpdateColumns(
		&dbmodel.Message{
			Id:     id,
			Status: string(status),
			Info:   strings.Join(info, "\t"),
		}, `status`, `info`); err != nil {
		if errors.Is(err, reform.ErrNoPK) {
			return fmt.Errorf(`%s %w %s`, id.String(), ErrNoMessage, err.Error())
		}
		if errors.Is(err, reform.ErrNoRows) {
			return fmt.Errorf(`%s %w %s`, id.String(), ErrNoMessage, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, id.String(), ErrDbUpdate, err.Error())
	}
	return nil
}

func (t *MessageRepository) strToStatus(status string) (model.Status, error) {
	switch status {
	case `new`:
		return model.StatusNew, nil
	case `processed`:
		return model.StatusProcessed, nil
	case `sent`:
		return model.StatusSent, nil
	case `error`:
		return model.StatusError, nil
	case `no_destinations`:
		return model.StatusNoDestinations, nil
	}
	return model.StatusError, fmt.Errorf(`%w %s`, ErrUnknownStatus, status)
}

func (t *MessageRepository) statusToSrt(status model.Status) string {
	switch status {
	case model.StatusNew:
		return `new`
	case model.StatusProcessed:
		return `processed`
	case model.StatusSent:
		return `sent`
	case model.StatusError:
		return `error`
	case model.StatusNoDestinations:
		return `no_destinations`
	}
	return fmt.Sprintf(`unknown (%s)`, status)
}
