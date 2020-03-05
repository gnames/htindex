package htindex_test

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
			Expect(hti.JobsNum).To(Equal(4))
			Expect(hti.OutputPath).To(Equal(testOutput))
		})
	})

	Describe("Run", func() {
		It("finds names in HathiTrust directory", func() {
			stdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.Run()).To(Succeed())
			data := getTestData(hti.OutputPath)
			// Issue #17 repetition of the same occurence many times in results
			hasRepetitions, err := hasRepetitions(hti)
			Expect(err).To(BeNil())
			Expect(hasRepetitions).To(BeFalse())
			Expect(hasWordsAround(data)).To(BeTrue())
			os.Stdout = stdout
		})

		// Issue #12
		It("reports about volumes that have no standard pages", func() {
			stdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.Run()).To(Succeed())
			errs, err := readErrors(hti.OutputPath)
			Expect(err).To(BeNil())

			res, ok := errs["yale.39002007302079"]
			Expect(ok).To(BeTrue())
			Expect(res.msg).To(Equal("non-standard naming for 76 pages"))
			os.Stdout = stdout
		})

		It("reports about volumes that have no pages", func() {
			stdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			hti, _ := NewHTindex(initOpts()...)
			Expect(hti.Run()).To(Succeed())
			errs, err := readErrors(hti.OutputPath)
			Expect(err).To(BeNil())

			res, ok := errs["yale.empty"]
			Expect(ok).To(BeTrue())
			Expect(res.msg).To(Equal("no pages detected"))
			os.Stdout = stdout
		})

		Measure("Going through titles fast enough", func(b Benchmarker) {
			stdout := os.Stdout
			os.Stdout, _ = os.Open(os.DevNull)
			runtime := b.Time("runtime", func() {
				hti, _ := NewHTindex(initOpts()...)
				Expect(hti.Run()).To(Succeed())
				os.Stdout = stdout
			})
			Expect(runtime.Seconds()).To(BeNumerically("<", 5.0), "took too long to run")
		}, 3)
	})
})

type testField int

const (
	timeStampF testField = iota
	idF
	pageIDF
	verbatimF
	nameStringF
	offsetStartF
	offsetEndF
	wordsBeforeF
	wordsAfterF
	oddsF
	kindF
)

type testData struct {
	TimeStamp   string
	ID          string
	PageID      string
	Verbatim    string
	NameString  string
	OffsetStart string
	OffsetEnd   string
	WordsBefore string
	WordsAfter  string
	Odds        string
	Kind        string
}

type htiError struct {
	ts      string
	titleID string
	pageID  string
	msg     string
}

func getTestData(path string) []testData {
	var res []testData
	resPath := filepath.Join(path, "results.csv")
	f, err := os.Open(resPath)
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
		datum := testData{
			TimeStamp:   v[timeStampF],
			ID:          v[idF],
			PageID:      v[pageIDF],
			Verbatim:    v[verbatimF],
			NameString:  v[nameStringF],
			OffsetStart: v[offsetStartF],
			OffsetEnd:   v[offsetEndF],
			WordsBefore: v[wordsBeforeF],
			WordsAfter:  v[wordsAfterF],
			Odds:        v[oddsF],
			Kind:        v[kindF],
		}
		res = append(res, datum)
	}

	return res
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
	root, err := filepath.Abs("./testdata")
	Expect(err).ToNot(HaveOccurred())
	input, err := filepath.Abs("./testdata/input_paths_small.txt")
	Expect(err).To(BeNil())
	opts := []Option{
		OptJobs(4),
		OptWordsAround(5),
		OptOutput("/tmp/htindex-test"),
		OptRoot(root),
		OptInput(input),
	}
	return opts
}

func hasRepetitions(hti *HTindex) (bool, error) {
	path := filepath.Join(hti.OutputPath, "results.csv")
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	ls, err := r.ReadAll()
	if err != nil {
		return false, err
	}
	if len(ls) < 200 {
		return false, fmt.Errorf("Result is too small: %d rows", len(ls))
	}

	indexID, indexPageID, indexOffsetStart, err := fieldIndices(ls[0])
	if err != nil {
		return false, err
	}
	data := make(map[string]struct{})
	for _, v := range ls[1:] {
		indices := []string{v[indexID], v[indexPageID], v[indexOffsetStart]}
		key := strings.Join(indices, "|")
		if _, ok := data[key]; ok {
			return true, nil
		} else {
			data[key] = struct{}{}
		}

	}
	return false, nil
}

func hasWordsAround(data []testData) bool {
	for _, d := range data {
		if len(d.WordsBefore) > 0 && len(d.WordsAfter) > 0 &&
			strings.Contains(d.WordsBefore, "|") &&
			strings.Contains(d.WordsAfter, "|") {
			return true
		}
	}
	return false
}

func fieldIndices(h []string) (int, int, int, error) {
	var titleID, pageID, offsetStart int
	for i, v := range h {
		switch v {
		case "ID":
			titleID = i
		case "PageID":
			pageID = i
		case "OffsetStart":
			offsetStart = i
		default:
			continue
		}
	}
	if titleID == 0 || pageID == 0 || offsetStart == 0 {
		err := fmt.Errorf(
			"Some indices are 0: titleID: %d, pageID: %d, offsetStart: %d",
			titleID, pageID, offsetStart,
		)
		return 0, 0, 0, err
	}
	return titleID, pageID, offsetStart, nil
}
