package htindex

import (
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Htindex", func() {
	Describe("NewHTindex", func() {
		It("creates an instance of HTindex", func() {
			hti, _ := NewHTindex()
			Expect(hti.jobsNum).To(Equal(runtime.NumCPU()))
			Expect(hti.progressNum).To(Equal(0))
		})

		It("can read options", func() {
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.jobsNum).To(Equal(15))
			Expect(hti.outputPath).To(Equal("/tmp/htindex-test"))
		})
	})

	Describe("Run", func() {
		It("finds names in HathiTrust directory", func() {
			stdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.Run()).To(Succeed())
			os.Stdout = stdout
		})
	})

})

func initOpts() []Option {
	root, err := filepath.Abs("./test-data")
	Expect(err).ToNot(HaveOccurred())
	input, err := filepath.Abs("./test-data/input_paths_small.txt")
	Expect(err).ToNot(HaveOccurred())
	opts := []Option{
		OptJobs(15),
		OptOutput("/tmp/htindex-test"),
		OptRoot(root),
		OptInput(input),
	}

	return opts
}
