package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
)

type Cmd struct {
	cli *cli.App
}

func (c *Cmd) Add(cmd *cli.Command) {
	c.cli.Commands = append(c.cli.Commands, cmd)
}

// Run is the entry point to the cli app. Parses the arguments slice and routes
// to the proper flag/args combination
func (c *Cmd) Run(ctx context.Context, arguments []string) error {
	return c.cli.RunContext(ctx, arguments)
}

func NewCmd() *Cmd {
	var cmd = &Cmd{
		cli: &cli.App{
			Flags: []cli.Flag{
				&cli.StringFlag{
					//@TODO: use HCL
					Name:     `config`,
					Aliases:  []string{`c`},
					Usage:    `path to config file`,
					Required: true,
				},
			},
			Commands: make([]*cli.Command, 0),
		},
	}
	return cmd
}
