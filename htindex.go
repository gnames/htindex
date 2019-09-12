package htindex

import (
	"fmt"
	"os"

	"github.com/gnames/gnfinder/dict"
)

// HTindex detects occurances of scientific names in Hathi Trust data.
type HTindex struct {
	// rootPrefix is concatenated with paths given in input file to get
	// complete path to HathiTrust files.
	rootPrefix string
	// inputPath gives path to file with input data.
	inputPath string
	// outputPath gives path to a directory to keep output data.
	outputPath string
	// jobsNum sets number of jobs/workers to run.
	jobsNum int
	// dict contains shared dictionary for name finding.
	dict *dict.Dictionary
	// reportNum determines how many titles should be processed for
	// a progress report.
	reportNum int
}

type Option func(h *HTindex)

func OptJobs(i int) Option {
	return func(h *HTindex) {
		h.jobsNum = i
	}
}

func OptReportNum(i int) Option {
	return func(h *HTindex) {
		h.reportNum = i
	}

}

func OptRoot(s string) Option {
	return func(h *HTindex) {
		h.rootPrefix = s
	}
}

func OptInput(s string) Option {
	return func(h *HTindex) {
		h.inputPath = s
	}
}

func OptOutput(s string) Option {
	return func(h *HTindex) {
		h.outputPath = s
	}
}

func NewHTindex(opts ...Option) (*HTindex, error) {
	hti := &HTindex{dict: dict.LoadDictionary(), reportNum: 0}
	for _, opt := range opts {
		opt(hti)
	}
	err := hti.setOutputDir()
	return hti, err
}

func (hti *HTindex) setOutputDir() error {
	path, err := os.Stat(hti.outputPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(hti.outputPath, 0755)
	}
	if path.Mode().IsRegular() {
		return fmt.Errorf("'%s' is a file, not a directory", hti.outputPath)
	}
	return nil
}
