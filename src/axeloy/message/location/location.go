package location

import "github.com/dmalykh/axeloy/axeloy/profile"

type Location interface {
	GetWay() string
	GetProfile() profile.Profile
}
