package htindex

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
}

type Option func(h *HTindex)

func OptJobs(i int) Option {
	return func(h *HTindex) {
		h.jobsNum = i
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

func NewHTindex(opts ...Option) *HTindex {
	hti := &HTindex{}
	for _, opt := range opts {
		opt(hti)
	}
	return hti
}
