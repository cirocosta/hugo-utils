package hugo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"

	"github.com/cirocosta/hugo-utils/hugo"
)

var _ = Describe("Pages", func() {
	Describe("GatherPages", func() {
		Context("with content directory having md pages w/ fm", func() {
			var (
				err   error
				pages []*hugo.Page
			)

			BeforeEach(func() {
				pages, err = hugo.GatherPages("testdata/content")
			})

			It("succeeds", func() {
				Expect(err).To(Succeed())
			})

			It("has path stored", func() {
				Expect(pages[0].Path).To(Equal("testdata/content/page1.md"))
				Expect(pages[1].Path).To(Equal("testdata/content/page2.md"))
			})

			It("has has frontmatter parsed", func() {
				Expect(pages[0].Title).To(Equal("page1"))
				Expect(pages[0].Date.Day()).To(Equal(2))
				Expect(len(pages[0].Tags)).To(Equal(2))
			})
		})
	})

	Describe("DiscoverMarkdownPaths", func() {
		Context("with empty root", func() {
			It("fails", func() {
				_, err := hugo.DiscoverMarkdownPaths("")
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with populated dir", func() {
			It("succeeds", func() {
				paths, err := hugo.DiscoverMarkdownPaths("testdata/content")

				Expect(err).To(Succeed())
				Expect(len(paths)).To(Equal(2))
				Expect(paths).To(ContainElement("testdata/content/page1.md"))
				Expect(paths).To(ContainElement("testdata/content/page2.md"))
			})
		})
	})

	Describe("ParsePageFile", func() {
		var err error

		Context("with empty path", func() {
			It("fails", func() {
				_, err = hugo.ParsePageFile("")
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with inexistent path", func() {
			It("fails", func() {
				_, err = hugo.ParsePageFile("/inexistent/path")
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with directory", func() {
			var tempDir string

			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "")
				Expect(err).To(Succeed())
			})

			AfterEach(func() {
				os.RemoveAll(tempDir)
			})

			It("fails", func() {
				_, err = hugo.ParsePageFile(tempDir)
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with an existing file", func() {
			var (
				filePath string
				page     *hugo.Page
			)

			Context("not having front matter", func() {
				It("fails", func() {
					filePath = "testdata/without-fm.md"
					page, err = hugo.ParsePageFile(filePath)
					Expect(err).ToNot(Succeed())
				})
			})

			Context("having front matter", func() {
				BeforeEach(func() {
					filePath = "testdata/page1.md"
					page, err = hugo.ParsePageFile(filePath)
				})

				It("succeeds", func() {
					Expect(err).To(Succeed())
				})

				It("has path properly set", func() {
					Expect(page.Path).To(Equal(filePath))
				})

				It("has front matter parsed", func() {
					Expect(page.Title).To(Equal("my thing"))
				})
			})
		})
	})
})
