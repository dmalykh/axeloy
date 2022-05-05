package cmd

import (
	"context"
	"github.com/urfave/cli/v2"
	"os"
)

type Cmd struct {
	cli *cli.App
}

func (c *Cmd) Add(cmd *cli.Command) {
	c.cli.Commands = append(c.cli.Commands, cmd)
}

func (c *Cmd) Run(ctx context.Context) error {
	return c.cli.Run(os.Args)
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
