package commands

import (
	_ "github.com/cirocosta/hugo-utils/hugo"
	"gopkg.in/urfave/cli.v1"
)

var Tags = cli.Command{
	Name:      "tags",
	Usage:     `manipulates tagged content.`,
	ArgsUsage: "[filename]",
	Action:    tagsAction,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "no-tags",
			Usage: "lists content without tags",
		},
	},
}

func tagsAction(c *cli.Context) (err error) {
	return
}
