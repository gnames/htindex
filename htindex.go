package htindex

import (
	"fmt"
	"os"
	"runtime"

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
	// progressNum determines how many titles should be processed for
	// a progress report.
	progressNum int
}

// Option sets the time for all options received during creation of new instance
// of HTindex object.
type Option func(h *HTindex)

// OptJobs sets number of jobs/workers to run duing execution.
func OptJobs(i int) Option {
	return func(h *HTindex) {
		h.jobsNum = i
	}
}

// OptProgressNum sets how often to printout a line about the progress. When it
// is set to 1 report line appears after processing every title, and if it is 10
// progress is shows after every 10th title.
func OptProgressNum(i int) Option {
	return func(h *HTindex) {
		h.progressNum = i
	}

}

// OptRoot sets the prefix of the path to zipped titles. It wil be concatenated
// with a path provided in the input file to receive complete absolute path.
func OptRoot(s string) Option {
	return func(h *HTindex) {
		h.rootPrefix = s
	}
}

// OptIntput is an absolute path to input data file. each line of such file
// displays path to zipped file of a title.
func OptInput(s string) Option {
	return func(h *HTindex) {
		h.inputPath = s
	}
}

// OptOutput is an absolute path to a directory where results will be written.
// If such directory does not exist already, it will be created during
// initialization of HTindex instance.
func OptOutput(s string) Option {
	return func(h *HTindex) {
		h.outputPath = s
	}
}

// NewHTindex creates HTindex instance with several defaults. If
// a some options are provided, they will override default settings.
func NewHTindex(opts ...Option) (*HTindex, error) {

	hti := &HTindex{
		dict:        dict.LoadDictionary(),
		progressNum: 0,
		jobsNum:     runtime.NumCPU(),
	}
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
