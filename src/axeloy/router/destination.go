package router

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
)

//Destination is profile with ways which can sent messages. One destination for one profile.
type Destination interface {
	GetWays(ctx context.Context) []way.Sender
	GetProfile(ctx context.Context) profile.Profile
	SetWays(ctx context.Context, senders ...way.Sender) error
	SetProfile(ctx context.Context, p profile.Profile) error
}
