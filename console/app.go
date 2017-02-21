package console

import "github.com/urfave/cli"

// Init console application
func InitApp() *cli.App {
	app := cli.NewApp()
	app.Name = "crawler"
	app.Usage = "analyser site"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run crawler full site",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "url",
					Value: "",
					Usage: "start url for analyse",
				},
			},
			Action: Parse,
		},
	}

	return app
}
