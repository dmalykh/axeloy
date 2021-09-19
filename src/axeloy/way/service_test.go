package way

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/way/repository"
	"reflect"
	"testing"
)

func TestWayService_GetAvailableListeners(t *testing.T) {
	type fields struct {
		wayRepository repository.WayRepository
		drivers       drivers
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Listener
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WayService{
				wayRepository: tt.fields.wayRepository,
				drivers:       tt.fields.drivers,
			}
			got, err := w.GetAvailableListeners(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvailableListeners() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAvailableListeners() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWayService_GetSenderByName(t *testing.T) {
	type fields struct {
		wayRepository repository.WayRepository
		drivers       drivers
	}
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Sender
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WayService{
				wayRepository: tt.fields.wayRepository,
				drivers:       tt.fields.drivers,
			}
			got, err := w.GetSenderByName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSenderByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSenderByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWayService_load(t *testing.T) {
	type fields struct {
		wayRepository repository.WayRepository
		drivers       drivers
	}
	type args struct {
		ctx    context.Context
		config Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WayService{
				wayRepository: tt.fields.wayRepository,
				drivers:       tt.fields.drivers,
			}
			if err := w.load(tt.args.ctx, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
