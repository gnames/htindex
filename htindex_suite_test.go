package htindex_test

import (
	"io/ioutil"
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testOutput = "/tmp/htindex-test"
)

func TestHtindex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Htindex Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
