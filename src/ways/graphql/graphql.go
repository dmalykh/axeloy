package graphql

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
)

type GraphQl struct {
	config config
}

type config struct {
	port string
}

func (q *GraphQl) Init(ctx context.Context, loader driver.ConfigLoader) error {
	if err := loader(&q.config); err != nil {
		return err
	}
	return nil
}
