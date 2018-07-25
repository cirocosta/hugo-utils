package commands

import (
	"io/ioutil"
	"os"

	"github.com/cirocosta/hugo-utils/hugo"
	"github.com/imdario/mergo"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
)

var Update = cli.Command{
	Name:  "update",
	Usage: "updates the frontmatter of a page.",
	Description: `The 'update' command takes care of updating the frontmatter
   of a given content page (e.g., /content/blog/mypost.md).

   Taking a desired update in the form of 'yaml', it parses the
   content page and applies to it the merge between the original
   frontmatter and the updated frontmatter.

   When no updated yaml is passed, the default frontmatter is
   applied (e.g., a post without 'tags' would now have 'tags: []').

EXAMPLES:

   Update the contents of page1.md with the defaults of the FrontMatter
   object from './hugo':

     cat ./page1.md
       ---
       title: 'page1'
       ---
       body

     hugo-utils update --filepath ./page1.md
     cat ./page1.md
       ---
       title: page1
       description: ""
       slug: ""
       image: ""
       date: 0001-01-01T00:00:00Z
       lastmod: 0001-01-01T00:00:00Z
       draft: false
       tags: []
       categories: []
       keywords: []
       ---
       body

   Update the contents of page1.md with the defaults of the FrontMatter
   object from './hugo' merged with a custom set of fields that we defined
   with a 'yaml' provided in the positional arguments:

     hugo-utils update --filepath ./page1.md 'tags: ["tag3"]'
     cat ./page1.md
       ---
       title: page1
       ... (other fields)
       tags:
       - tag1
       - tag2
       ---
       body
`,
	Action:    updateAction,
	ArgsUsage: "[yaml]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "filepath",
			Usage: "path to the page file",
		},
	},
}

func updateAction(c *cli.Context) (err error) {
	var (
		pageFilepath = c.String("filepath")
		yamlSrc      = c.Args().First()
		updateFm     = &hugo.FrontMatter{}
		page         *hugo.Page
		file         *os.File
		tempFile     *os.File
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

	page, err = hugo.ParsePage(file)
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}

	if yamlSrc != "" {
		err = yaml.Unmarshal([]byte(yamlSrc), updateFm)
		if err != nil {
			err = cli.NewExitError(err, 1)
			return
		}

		err = mergo.Merge(&page.FrontMatter, updateFm, mergo.WithOverride)
		if err != nil {
			err = cli.NewExitError(err, 1)
			return
		}
	}

	tempFile, err = ioutil.TempFile("", "")
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}
	defer tempFile.Close()

	err = page.Write(tempFile)
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}

	err = tempFile.Close()
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}

	err = os.Rename(tempFile.Name(), pageFilepath)
	if err != nil {
		err = cli.NewExitError(err, 1)
		return
	}

	return
}
