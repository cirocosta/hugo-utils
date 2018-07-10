package hugo

import (
	"os"

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

// ParsePage parses a single page given a filepath.
func ParsePage(path string) (page *Page, err error) {
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

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)

	fm, _, err := m.Parse(file)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse %s front matter")
		return
	}

	page = &Page{
		Path:        path,
		FrontMatter: fm,
	}

	return
}

// GatherPages lists all the pages found after a given root.
func GatherPages(root string) (pages []*Page, err error) {
	return
}
