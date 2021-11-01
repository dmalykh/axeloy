package location

import "github.com/dmalykh/axeloy/axeloy/profile"

type Location interface {
	GetWays() []string
	GetProfile() profile.Profile
}
