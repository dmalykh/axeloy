package model

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
)

type RouteDestination struct {
	Profile profile.Profile
	Ways    []way.Sender
}

func (r *RouteDestination) GetWays(ctx context.Context) []way.Sender {
	return r.Ways
}

func (r *RouteDestination) GetProfile(ctx context.Context) profile.Profile {
	return r.Profile
}

func (r *RouteDestination) SetWays(ctx context.Context, senders ...way.Sender) error {
	r.Ways = senders
	return nil
}

func (r *RouteDestination) SetProfile(ctx context.Context, p profile.Profile) error {
	r.Profile = p
	return nil
}
