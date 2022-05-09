package main

import (
	"context"
	"github.com/dmalykh/axeloy/cmd"
	"github.com/dmalykh/axeloy/cmd/app"
	"github.com/dmalykh/axeloy/cmd/atlas"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	var cli = cmd.NewCmd()
	cli.Add(app.Command(ctx))
	cli.Add(atlas.Command(ctx))

	if err := cli.Run(ctx, os.Args); err != nil {
		cancel()
		log.Fatal(err)
	}
}
