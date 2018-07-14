package commands

import (
	"fmt"
	"os"
	"path"
	"text/tabwriter"
	"text/template"

	"github.com/cirocosta/hugo-utils/hugo"
	"gopkg.in/urfave/cli.v1"
)

var List = cli.Command{
	Name: "list",
	Usage: `lists all content under a given path.

   The 'list' command iterates over each content file (*.md)
   found under a given root directory (--directory), then prints
   to 'stdout' a description of each.

   The default formatting displays the following attributes for
   each page: title, keywords, tags, categories, slug, date.

   A custom format can also be specified following Go template
   rules. In this case, the render state contains:
   - {{ . }}: the current page in the page traversal; and
   - {{ .Pages }}: the list of all pages found.

EXAMPLES:

   Display every property of the pages under a given
   section that lives under "./content/blog" using the default
   formatting:

     hugo-utils \
       --directory=./content/blog

   Display the text of every page in a given section
   that lives under "./content/blog" and their keywords:

     hugo-utils \
       --directory=./content/blog \
       '{{ .Title }} - {{ .Keywords }}'

   Display the path to the files that don't have keywords
   specified:

     hugo-utils \
       --directory=./content/blog \
       '{{ if eq (len .Keywords) 0 }} {{ .Path }} {{ end }}'
`,
	ArgsUsage: "[format]",
	Action:    listAction,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "directory",
			Usage: "path to the directory where contents exist (.md)",
		},
	},
}

type renderState struct {
	*hugo.Page
	Pages []*hugo.Page
}

func listAction(c *cli.Context) (err error) {
	var (
		root   = c.String("directory")
		format = c.Args().First()
	)

	if root == "" {
		cli.ShowCommandHelp(c, "list")
		err = cli.NewExitError("a root path must be specified", 1)
		return
	}

	pages, err := hugo.GatherPages(root)
	if err != nil {
		cli.ShowCommandHelp(c, "list")
		err = cli.NewExitError(err, 1)
		return
	}

	if format == "" {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
		for _, page := range pages {
			fmt.Fprintf(w, "%s\t%s\n", "title", page.Title)
			fmt.Fprintf(w, "%s\t%v\n", "file", path.Base(page.Path))
			fmt.Fprintf(w, "%s\t%v\n", "slug", page.Slug)
			fmt.Fprintf(w, "%s\t%v\n", "date", page.Date.Format("Jan 2, 2006"))
			fmt.Fprintf(w, "%s\t%v\n", "last-mod", page.LastMod.Format("Jan 2, 2006"))
			fmt.Fprintf(w, "%s\t%v\n", "keywords", page.Keywords)
			fmt.Fprintf(w, "%s\t%v\n", "tags", page.Tags)
			fmt.Fprintf(w, "%s\t%v\n", "categories", page.Categories)
			fmt.Fprintf(w, "\n")
		}
		w.Flush()
		return
	}

	t, err := template.New("list-format").Parse(format)
	if err != nil {
		cli.ShowCommandHelp(c, "list")
		err = cli.NewExitError(err, 1)
		return
	}

	for _, page := range pages {
		err = t.Execute(os.Stdout, &renderState{page, pages})
		if err != nil {
			err = cli.NewExitError(err, 1)
		}

		fmt.Fprintln(os.Stdout, "")
	}

	return
}
