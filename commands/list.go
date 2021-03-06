package commands

import (
	"fmt"
	"os"
	"path"
	"sort"
	"text/tabwriter"
	"text/template"

	"github.com/cirocosta/hugo-utils/hugo"
	"gopkg.in/urfave/cli.v1"
)

var List = cli.Command{
	Name:  "list",
	Usage: "lists all content under a given path.",
	Description: `The 'list' command iterates over each content file (*.md)
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
		cli.StringFlag{
			Name:  "type",
			Usage: "content type to list entries by (pages|tags|categories)",
			Value: "pages",
		},
		cli.StringFlag{
			Name:  "sort",
			Usage: "thing to sort by (title|date|lastmod)",
			Value: "lastmod",
		},
		cli.BoolFlag{
			Name:  "draft",
			Usage: "only show drafts",
		},
	},
}

type renderState struct {
	*hugo.Page
	Pages []*hugo.Page
}

func showPagesList(c *cli.Context, pages []*hugo.Page) {
	var (
		format = c.String("format")
		draft  = c.Bool("draft")
	)

	if format == "" {
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
		for _, page := range pages {
			if draft {
				if !page.Draft {
					continue
				}
			}

			fmt.Fprintf(w, "%s\t%s\n", "title", page.Title)
			fmt.Fprintf(w, "%s\t%v\n", "file", path.Base(page.Path))
			fmt.Fprintf(w, "%s\t%v\n", "slug", page.Slug)
			fmt.Fprintf(w, "%s\t%v\n", "date", page.Date.Format("Jan 2, 2006"))
			fmt.Fprintf(w, "%s\t%v\n", "last-mod", page.LastMod.Format("Jan 2, 2006"))
			fmt.Fprintf(w, "%s\t%v\n", "keywords", page.Keywords)
			fmt.Fprintf(w, "%s\t%v\n", "tags", page.Tags)
			fmt.Fprintf(w, "%s\t%v\n", "draft", page.Draft)
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
}

func showTagsMap(c *cli.Context, pages []*hugo.Page) {
	var tagsMapping = map[string][]*hugo.Page{}

	for _, page := range pages {
		if len(page.Tags) == 0 {
			continue
		}

		for _, tag := range page.Tags {
			mapping, ok := tagsMapping[tag]
			if !ok {
				tagsMapping[tag] = []*hugo.Page{page}
				continue
			}

			mapping = append(mapping, page)
			tagsMapping[tag] = mapping
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 4, ' ', 0)
	for tag, tagPages := range tagsMapping {
		fmt.Fprintf(w, "%s\n", tag)
		for _, page := range tagPages {
			fmt.Fprintf(w, "\t%s\t(%s)\n", page.Title, path.Base(page.Path))
		}
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\n")
	}

	w.Flush()
	return

}

func listAction(c *cli.Context) (err error) {
	var (
		root     = c.String("directory")
		listType = c.String("type")
		sortBy   = c.String("sort")
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

	if sortBy != "" {
		switch sortBy {
		case "title":
			sort.Slice(pages, func(i, j int) bool {
				return pages[i].Title < pages[j].Title
			})
		case "date":
			sort.Slice(pages, func(i, j int) bool {
				return pages[i].Date.Before(pages[j].Date)
			})
		case "lastmod":
			sort.Slice(pages, func(i, j int) bool {
				return pages[i].LastMod.Before(pages[j].LastMod)
			})
		default:
			cli.ShowCommandHelp(c, "list")
			err = cli.NewExitError("unknown sort type", 1)
			return
		}
	}

	switch listType {
	case "tags":
		showTagsMap(c, pages)
	case "pages":
		showPagesList(c, pages)
	default:
		cli.ShowCommandHelp(c, "list")
		err = cli.NewExitError("unknown list type "+listType, 1)
		return
	}

	return
}
