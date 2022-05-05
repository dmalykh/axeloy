package atlas

import (
	"context"
	"fmt"
	configuration "github.com/dmalykh/axeloy/config"
	"github.com/urfave/cli/v2"
)

func Command(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  `atlas`,
		Usage: `run atlas migrations tool, look for more info at https://atlasgo.io/`,
		Subcommands: []*cli.Command{
			{
				Name: `inspect`,
				Action: func(c *cli.Context) error {
					//Parse config
					config, err := configuration.Load(c.String(`config`))
					if err != nil {
						return cli.Exit(fmt.Sprintf(`open config error %s`, err.Error()), 8)
					}
					// Init atlas
					a, err := NewAtlas(ctx, config.Db.Driver, config.Db.Dsn)
					if err != nil {
						return cli.Exit(fmt.Sprintf(`open create atlas %s`, err.Error()), 8)
					}
					a.Inspect(ctx)

					return nil
				},
			},
		},
	}
}
