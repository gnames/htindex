package htindex

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gnames/gnfinder/output"
)

// detectedName holds information about a name-string returned by a
// name-finder.
type detectedName struct {
	pageID      string
	verbatim    string
	nameString  string
	offsetStart int
	offsetEnd   int
	wordsBefore string
	wordsAfter  string
	annotNomen  string
	odds        float64
	kind        string
	timestamp   string
}

// outputError outputs errors arrived from the name-finding process.
func (hti *HTindex) outputError(errCh <-chan *htiError, wgOut *sync.WaitGroup) {
	f, err := os.Create(filepath.Join(hti.OutputPath, "errors.csv"))
	defer wgOut.Done()
	if err != nil {
		log.Fatal(err)
	}
	ef := csv.NewWriter(f)
	_ = ef.Write([]string{"TimeStamp", "TitleID", "PageID", "Error"})
	defer f.Close()
	defer ef.Flush()
	for e := range errCh {
		_ = ef.Write([]string{e.ts, e.titleID, e.pageID, e.msg})
	}
}

// outputResults outputs data about found names.
func (hti *HTindex) outputResult(outCh <-chan *title, wgOut *sync.WaitGroup) {
	defer wgOut.Done()
	count := 0
	ts := time.Now()

	f, err := os.Create(filepath.Join(hti.OutputPath, "results.csv"))
	if err != nil {
		log.Fatal(err)
	}
	titles, err := os.Create(filepath.Join(hti.OutputPath, "titles.csv"))
	if err != nil {
		log.Fatal(err)
	}

	of := csv.NewWriter(f)
	tf := csv.NewWriter(titles)
	_ = of.Write([]string{
		"TimeStamp", "ID", "PageID", "Verbatim", "WordsBefore", "NameString",
		"WordsAfter", "AnnotNomen", "OffsetStart", "OffsetEnd", "Odds", "Kind",
	})
	_ = tf.Write([]string{"ID", "SHA256", "Path", "PagesNumber", "BadPagesNumber", "NamesOccurences"})

	defer f.Close()
	defer titles.Close()
	defer of.Flush()
	defer tf.Flush()

	for t := range outCh {
		_ = tf.Write([]string{
			t.id, t.sha256, t.path, strconv.Itoa(len(t.pages)),
			strconv.Itoa(t.pagesNumBadNames), strconv.Itoa(t.namesNum),
		})

		count++
		if hti.ProgressNum > 0 && count%hti.ProgressNum == 0 {
			rate := float64(count) / (time.Since(ts).Minutes())
			log.Printf("Processing %dth title. Rate %0.2f titles/min\n", count, rate)
		}
		if t.namesNum == 0 {
			continue
		}
		for _, p := range t.pages {
			for _, name := range p.res.Names {
				n := newDetectedName(p, name)
				out := []string{
					n.timestamp, t.id, n.pageID, n.verbatim, n.wordsBefore,
					n.nameString, n.wordsAfter, n.annotNomen,
					strconv.Itoa(n.offsetStart), strconv.Itoa(n.offsetEnd),
					strconv.Itoa(int(n.odds)), n.kind,
				}
				_ = of.Write(out)

				if err := of.Error(); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// ts generates a converted to a string timestamp in nanoseconds from epoch.
func ts() string {
	t := time.Now()
	return strconv.FormatInt(t.UnixNano(), 10)
}

// newDetectedName processes output from name-finding to prepare it for
// htindex output.
func newDetectedName(p page, n output.Name) detectedName {
	dn := detectedName{
		pageID:      p.id,
		verbatim:    n.Verbatim,
		nameString:  n.Name,
		offsetStart: n.OffsetStart,
		offsetEnd:   n.OffsetEnd,
		wordsBefore: strings.Join(n.WordsBefore, "|"),
		wordsAfter:  strings.Join(n.WordsAfter, "|"),
		annotNomen:  n.AnnotNomen,
		odds:        n.Odds,
		kind:        n.Type,
		timestamp:   ts(),
	}
	return dn
}
