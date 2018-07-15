package hugo

import (
	"time"
)

// FrontMatter corresponds to the parsed front
// matter of the page.
type FrontMatter struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Slug        string    `yaml:"slug"`
	Image       string    `yaml:"image"`
	Date        time.Time `yaml:"date"`
	LastMod     time.Time `yaml:"lastmod"`
	Tags        []string  `yaml:"tags"`
	Categories  []string  `yaml:"categories"`
	Keywords    []string  `yaml:"keywords"`
	Draft       bool      `yaml:"draft"`
}
