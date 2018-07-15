package main

import (
	"os"

	"github.com/cirocosta/hugo-utils/commands"
	"gopkg.in/urfave/cli.v1"
)

var (
	version = "dev"
	commit  = "HEAD"
)

func main() {
	app := cli.NewApp()

	app.Version = version + " - " + commit
	app.Usage = "missing hugo tools"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Ciro da Silva da Costa",
			Email: "ciro.costa@liferay.com",
		},
	}
	app.Description = `hugo-utils provides missing features from the "hugo" CLI tool`
	app.Commands = []cli.Command{
		commands.List,
		commands.Update,
	}

	app.Run(os.Args)
}
