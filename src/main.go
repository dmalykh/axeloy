package main

import (
	"context"
	"github.com/dmalykh/axeloy/cmd"
	"github.com/dmalykh/axeloy/cmd/app"
	"github.com/dmalykh/axeloy/cmd/atlas"
	"log"
)

func main() {
	var ctx = context.Background()
	var cli = cmd.NewCmd()
	cli.Add(app.Command(ctx))
	cli.Add(atlas.Command(ctx))

	if err := cli.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
