package htindex_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gnames/htindex"
)

var _ = Describe("Htindex", func() {
	Describe("NewHTindex", func() {
		It("creates an instance of HTindex", func() {
			hti, _ := NewHTindex()
			Expect(hti.JobsNum).To(Equal(runtime.NumCPU()))
			Expect(hti.ProgressNum).To(Equal(0))
		})

		It("can read options", func() {
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.JobsNum).To(Equal(15))
			Expect(hti.OutputPath).To(Equal(testOutput))
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

		// Issue #12
		It("reports about volumes that do not have expected page names", func() {
			stdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.Run()).To(Succeed())
			errs, err := readErrors(hti.OutputPath)
			Expect(err).To(BeNil())

			res, ok := errs["yale.39002007302079"]
			Expect(ok).To(BeTrue())
			Expect(res.msg).To(Equal("no pages detected"))
			os.Stdout = stdout
		})
	})
})

type htiError struct {
	ts      string
	titleID string
	pageID  string
	msg     string
}

func readErrors(path string) (map[string]*htiError, error) {
	res := make(map[string]*htiError)
	errPath := filepath.Join(path, "errors.csv")
	f, err := os.Open(errPath)
	Expect(err).To(BeNil())
	reader := csv.NewReader(f)
	count := 0
	for {
		count++
		v, err := reader.Read()
		if err == io.EOF {
			break
		}
		if count == 1 {
			continue
		}
		Expect(err).To(BeNil())
		fmt.Println(v)
		res[v[1]] = &htiError{
			ts:      v[0],
			titleID: v[1],
			pageID:  v[2],
			msg:     v[3],
		}
	}
	Expect(count).To(BeNumerically(">", 0))
	return res, nil
}

func initOpts() []Option {
	root, err := filepath.Abs("./test-data")
	Expect(err).ToNot(HaveOccurred())
	input, err := filepath.Abs("./test-data/input_paths_small.txt")
	Expect(err).To(BeNil())
	opts := []Option{
		OptJobs(15),
		OptOutput("/tmp/htindex-test"),
		OptRoot(root),
		OptInput(input),
	}
	return opts
}
