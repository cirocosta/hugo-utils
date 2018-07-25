package hugo

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Page represents a content page.
type Page struct {
	// Path is path to the page file.
	Path string

	// FrontMatter corresponds to the parsed front
	// matter of the page.
	FrontMatter `yaml:"-,inline"`

	// Body contains the actual content of the page.
	Body []byte
}

var frontMatterDelim = []byte("---")

type (
	ParseState uint8
)

const (
	ParseStateStart ParseState = iota
	ParseStateDelimStart
	ParseStateFrontMatter
	ParseStateDelimEnd
	ParseStateBody
)

// SplitFrontMatterAndBody takes a given reader and then
// splits its content in two:
// - FrontMatter
// - Body
func SplitFrontMatterAndBody(r io.Reader) (frontMatter, body []byte, err error) {
	if r == nil {
		err = errors.Errorf(
			"a reader must be specified")
		return
	}

	var (
		text            []byte
		scanner         = bufio.NewScanner(r)
		delimetersFound = 0
		state           = ParseStateStart
	)

	body = make([]byte, 0)
	frontMatter = make([]byte, 0)

	for scanner.Scan() {
		text = scanner.Bytes()

		if delimetersFound < 2 {
			if bytes.Equal(frontMatterDelim, text) {
				delimetersFound++

				if delimetersFound == 1 {
					state = ParseStateFrontMatter
				} else {
					state = ParseStateDelimEnd
				}
			}
		}

		switch state {
		case ParseStateBody:
			body = append(body, text...)
			body = append(body, '\n')
		case ParseStateFrontMatter:
			frontMatter = append(frontMatter, text...)
			frontMatter = append(frontMatter, '\n')
		case ParseStateDelimStart:
			continue
		case ParseStateDelimEnd:
			state = ParseStateBody
		}
	}

	err = scanner.Err()
	if err != nil {
		err = errors.Wrapf(err,
			"failed while scanning page")
		return
	}

	return
}

// Write writes the contents of a page to a given destination
// writer.
func (p *Page) Write(w io.Writer) (err error) {
	if w == nil {
		err = errors.Errorf("writer msut not be nil")
		return
	}

	_, err = w.Write([]byte("---\n"))
	if err != nil {
		err = errors.Wrapf(err, "failed to write openning front matter separator")
		return
	}

	encoder := yaml.NewEncoder(w)
	err = encoder.Encode(&p.FrontMatter)
	if err != nil {
		err = errors.Wrapf(err, "failed to encode frontmatter to document")
		return
	}

	err = encoder.Close()
	if err != nil {
		err = errors.Wrapf(err, "failed to close frontmatter encoder")
		return
	}

	_, err = w.Write([]byte("---\n"))
	if err != nil {
		err = errors.Wrapf(err, "failed to write closing front matter separator")
		return
	}

	r := bytes.NewReader(p.Body)
	_, err = io.Copy(w, r)
	if err != nil {
		err = errors.Wrapf(err, "failed to copy body to writer")
		return
	}

	return
}

// ParsePage parses the page contents.
func ParsePage(r io.Reader) (page *Page, err error) {
	front, body, err := SplitFrontMatterAndBody(r)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to split frommatter and body from content page")
		return
	}

	page = &Page{Body: body}
	err = yaml.Unmarshal(front, &page.FrontMatter)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse front matter from reader")
		return
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

	_, err = os.Stat(root)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve info from root path %s", root)
		return
	}

	paths = make([]string, 0)
	walkFunc := func(path string, info os.FileInfo, walkErr error) (err error) {
		if err != nil {
			return
		}

		if info == nil {
			err = errors.Errorf("couldn't gather info from path %s", path)
			return
		}

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
