package smsru

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
	"github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strings"
)

type Config struct {
}

type SmsRu struct {
	url                string
	login              string
	password           string
	defaultCountryCode string
	strictValidation   bool
	params             params
}

type sms struct {
	phone string
	text  string
}

type params struct {
	text string
}

func (s *SmsRu) ValidateProfile(ctx context.Context, p profile.Profile) error {
	return validation.Validate(p,
		validation.Map(
			validation.Key("phone", validation.Required, IsPhoneNumber(s.defaultCountryCode, s.strictValidation)),
		),
	)
}

func (s *SmsRu) SetWayParams(params driver.Params) {
	s.params.text = params[`text`]
}

func (s *SmsRu) SetConfig(config driver.DriverConfig) {
	//@TODO configparser. driver.DriverConfig should be interface{}. config.Load(s, config)
}

func (s *SmsRu) Stop() error {
	panic("implement me")
}

func (s *SmsRu) Send(ctx context.Context, recipient profile.Profile, message message.Message) ([]string, error) {
	var sms sms
	if err := profile.Unmarshal(recipient, &sms); err != nil {

	}
	sms.text = strings.NewReplacer().Replace(s.)

	http.Post(s.url, ``, sms)
}

var Driver SmsRu
