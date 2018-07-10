package hugo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHugo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hugo Suite")
}
