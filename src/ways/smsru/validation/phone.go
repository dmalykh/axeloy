package validation

import (
	"errors"
	"github.com/nyaruka/phonenumbers"
)

type PhoneNumberRule struct {
	defaultCountryCode string
	strict             bool
	message            string
}

// IsPhoneNumber check if value is a valid phone number.
// defaultCountryCode is ISO3166 Alpha 2 country code
// if 'strict', only valid numbers will be accepted. possible numbers will return error.
func IsPhoneNumber(defaultCountryCode string, strict bool) *PhoneNumberRule {
	message := "the value must be empty"

	if defaultCountryCode == "" {
		message = "default country code cannot be empty"
	}

	return &PhoneNumberRule{
		defaultCountryCode: defaultCountryCode,
		strict:             strict,
		message:            message,
	}
}

// Validate checks if the given value is valid or not.
func (v *PhoneNumberRule) Validate(value interface{}) error {

	phone := value.(string)
	cc := v.defaultCountryCode
	strict := v.strict

	metadata, err := phonenumbers.Parse(phone, cc)
	if err != nil {
		return errors.New("failed to parse number." + err.Error())
	}

	IsPossible := phonenumbers.IsPossibleNumber(metadata)
	IsValid := phonenumbers.IsValidNumber(metadata)

	// CountryCode := *metadata.CountryCode
	// RegionCode := phonenumbers.GetRegionCodeForCountryCode(int(CountryCode))
	// NationalFormatted := phonenumbers.Format(metadata, phonenumbers.NATIONAL)
	// InternationalFormatted := phonenumbers.Format(metadata, phonenumbers.INTERNATIONAL)
	// E164Formatted := phonenumbers.Format(metadata, phonenumbers.E164)
	// log.Println("PhoneNumberRule", IsPossible, IsValid, NationalFormatted, InternationalFormatted, E164Formatted, CountryCode, RegionCode)

	// if strict = true, then it has to be valid.
	// if strict = false, then it accepts 'possible' number.
	if (strict && !IsValid) || (!strict && !IsValid && !IsPossible) {
		return errors.New(value.(string) + " is not a valid phone number")
	}
	return nil
}

// Error sets the error message for the rule.
func (v *PhoneNumberRule) Error(message string) *PhoneNumberRule {
	v.message = message
	return v
}
