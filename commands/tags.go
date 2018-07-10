package commands

import (
	"gopkg.in/urfave/cli.v1"
)

var Tags = cli.Command{
	Name:      "tags",
	Usage:     `manipulates tagged content.`,
	ArgsUsage: "[filename]",
	Action:    tagsAction,
	Flags:     []cli.Flag{},
}

func tagsAction(c *cli.Context) (err error) {
	return
}
