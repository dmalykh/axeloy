package model

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
	mocks "github.com/dmalykh/axeloy/mocks/axeloy/profile"
	mocks2 "github.com/dmalykh/axeloy/mocks/axeloy/way"
	"reflect"
	"testing"
)

func TestRouteDestination(t *testing.T) {
	type fields struct {
		Profile profile.Profile
		Ways    []way.Sender
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			fields: fields{
				Profile: &mocks.Profile{},
				Ways: []way.Sender{
					&mocks2.Sender{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()
			r := &RouteDestination{}

			if err := r.SetWays(ctx, tt.fields.Ways...); err != nil {
				t.Errorf("SetWays(%v) error %s", tt.fields.Ways, err.Error())
			}
			if err := r.SetProfile(ctx, tt.fields.Profile); err != nil {
				t.Errorf("SetProfile(%v) error %s", tt.fields.Profile, err.Error())
			}

			if got := r.GetProfile(ctx); !reflect.DeepEqual(got, tt.fields.Profile) {
				t.Errorf("GetProfile() = %v, want %v", got, tt.fields.Profile)
			}
			if got := r.GetWays(ctx); !reflect.DeepEqual(got, tt.fields.Ways) {
				t.Errorf("GetWays() = %v, want %v", got, tt.fields.Ways)
			}
		})
	}
}
