package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cirocosta/hugo-utils/hugo"
	"gopkg.in/urfave/cli.v1"
)

var List = cli.Command{
	Name:      "list",
	Usage:     `lists all content under a given path`,
	ArgsUsage: "[content directory]",
	Action:    listAction,
}

func listAction(c *cli.Context) (err error) {
	var (
		root = c.Args().First()
	)

	if root == "" {
		err = cli.NewExitError("a root path must be specified", 1)
		return
	}

	pages, err := hugo.GatherPages(root)
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	for _, page := range pages {
		fmt.Fprintf(w, "%s\t%s\n", "title", page.Title)
		fmt.Fprintf(w, "%s\t%v\n", "keywords", page.Keywords)
		fmt.Fprintf(w, "%s\t%v\n", "tags", page.Tags)
		fmt.Fprintf(w, "%s\t%v\n", "categories", page.Categories)
		fmt.Fprintf(w, "%s\t%v\n", "slug", page.Slug)
		fmt.Fprintf(w, "%s\t%v\n", "date", page.Date.Format("Jan 2, 2006"))
		fmt.Fprintf(w, "\n")
	}
	w.Flush()

	return
}
