package app

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
)

func Command(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  `run`,
		Usage: `run axeloy`,
		Action: func(c *cli.Context) error {
			// Get app
			var application = NewApp()
			// Open config
			config, err := application.Open(c.String(`config`))
			if err != nil {
				return cli.Exit(fmt.Sprintf(`open config error %s`, err.Error()), 8)
			}
			// Load app
			var ctx = application.WithShutdown(ctx)
			ax, err := application.Load(ctx, config)
			if err != nil {
				return cli.Exit(fmt.Sprintf(`config loading error %s`, err.Error()), 9)
			}
			// Run
			return ax.Run(ctx)
		},
	}
}
