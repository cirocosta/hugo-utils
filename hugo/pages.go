package hugo

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gernest/front"
	"github.com/pkg/errors"
)

// Page represents a content page.
type Page struct {
	// Path is path to the page file.
	Path string

	// FrontMatter corresponds to the parsed front
	// matter of the page.
	FrontMatter map[string]interface{}
}

// ParsePage parses the page contents.
func ParsePage(r io.Reader) (page *Page, err error) {
	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)

	fm, _, err := m.Parse(r)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse front matter from reader")
		return
	}

	page = &Page{
		FrontMatter: fm,
	}

	return
}

// ParsePageFile parses a single page given a filepath.
func ParsePageFile(path string) (page *Page, err error) {
	if path == "" {
		err = errors.Errorf("path must be non-empty")
		return
	}

	file, err := os.Open(path)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to open file %s",
			path)
		return
	}
	defer file.Close()

	page, err = ParsePage(file)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse %s page content", path)
		return
	}
	page.Path = path

	return
}

// DiscoverMarkdownPaths looks at the filesystem as indicated
// by a root path and searches for markdown files that might live there.
func DiscoverMarkdownPaths(root string) (paths []string, err error) {
	if root == "" {
		err = errors.Errorf("a root must be specified")
		return
	}

	paths = make([]string, 0)
	walkFunc := func(path string, info os.FileInfo, walkErr error) (err error) {
		if info.IsDir() {
			return
		}

		if filepath.Ext(path) != ".md" {
			return
		}

		paths = append(paths, path)
		return
	}

	err = filepath.Walk(root, walkFunc)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to search for markdown files under root %s",
			root)
		return
	}

	return
}

// GatherPages gathers all content pages under a given root path
// and parses their contents.
func GatherPages(root string) (pages []*Page, err error) {
	var (
		paths []string
		page  *Page
	)

	paths, err = DiscoverMarkdownPaths(root)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't find content under %s", root)
		return
	}

	for _, path := range paths {
		page, err = ParsePageFile(path)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to parse page %s", path)
			return
		}

		pages = append(pages, page)
	}

	return
}
