package htindex_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHtindex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Htindex Suite")
}
