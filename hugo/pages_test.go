package hugo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"io/ioutil"

	"github.com/cirocosta/hugo-utils/hugo"
)

var _ = Describe("Pages", func() {
	Describe("ParsePage", func() {
		var err error

		Context("with empty path", func() {
			It("fails", func() {
				_, err = hugo.ParsePage("")
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with inexistent path", func() {
			It("fails", func() {
				_, err = hugo.ParsePage("/inexistent/path")
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with directory", func() {
			var tempDir string

			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "")
				Expect(err).To(Succeed())
			})

			AfterEach(func () {
				os.RemoveAll(tempDir)
			})

			It("fails", func() {
				_, err = hugo.ParsePage(tempDir)
				Expect(err).ToNot(Succeed())
			})
		})

		Context("with an existing file", func () {
			var (
				filePath string
				page *hugo.Page
			)

			Context("not having front matter", func () {
				It("fails", func () {
					filePath = "testdata/without-fm.md"
					page, err = hugo.ParsePage(filePath)
					Expect(err).ToNot(Succeed())
				})
			})

			Context("having front matter", func () {
				BeforeEach(func () {
					filePath = "testdata/page1.md"
					page, err = hugo.ParsePage(filePath)
				})

				It("succeeds", func () {
					Expect(err).To(Succeed())
				})

				It("has path properly set", func () {
					Expect(page.Path).To(Equal(filePath))
				})

				It("has front matter parsed", func () {
					_, ok := page.FrontMatter["title"]
					Expect(ok).To(BeTrue())
				})
			})
		})
	})
})
