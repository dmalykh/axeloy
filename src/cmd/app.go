package cmd

import (
	"fmt"
	"github.com/dmalykh/axeloy/cmd/app"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Root() {
	var cmd = &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run axeloy",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     `config`,
						Aliases:  []string{`c`},
						Usage:    `path to config file`,
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					var a = new(app.App)
					var ctx = a.NewContext()
					ax, err := a.Load(ctx, c.String(`config`))
					if err != nil {
						return cli.Exit(fmt.Sprintf(`config lodaing error %s`, err.Error()), 9)
					}
					return ax.Run(ctx)
				},
			},
		},
	}

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
