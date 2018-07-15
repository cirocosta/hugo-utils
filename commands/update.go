package commands

import (
	"os"

	"github.com/cirocosta/hugo-utils/hugo"
	"gopkg.in/urfave/cli.v1"
)

var Update = cli.Command{
	Name:      "update",
	Usage:     `updates the frontmatter of a page.`,
	Action:    updateAction,
	ArgsUsage: "filepath",
}

func updateAction(c *cli.Context) (err error) {
	var (
		pageFilepath = c.Args().First()
		file         *os.File
	)

	if pageFilepath == "" {
		cli.ShowCommandHelp(c, "update")
		err = cli.NewExitError("a filepath must be specified", 1)
		return
	}

	file, err = os.Open(pageFilepath)
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}
	defer file.Close()

	_, err = hugo.ParsePage(file)
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}

	return
}
