package app

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy"
	"log"
)

func main() {
	var ax axeloy.Axeloy
	var ctx = context.Background()
	if err := ax.Run(ctx, &axeloy.Config{}); err != nil {
		log.Fatalln(err.Error())
	}
}
