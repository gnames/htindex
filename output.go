package htindex

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gnames/gnfinder/output"
)

// detectedName holds information about a name-string returned by a
// name-finder.
type detectedName struct {
	pageID       string
	verbatim     string
	nameString   string
	offsetStart  int
	offsetEnd    int
	endsNextPage int
	odds         float64
	kind         string
	timestamp    string
}

// outputError outputs errors arrived from the name-finding process.
func (hti *HTindex) outputError(errCh <-chan *htiError, wgOut *sync.WaitGroup) {
	f, err := os.Create(filepath.Join(hti.OutputPath, "errors.csv"))
	defer wgOut.Done()
	if err != nil {
		log.Fatal(err)
	}
	ef := csv.NewWriter(f)
	ef.Write([]string{"TimeStamp", "TitleID", "PageID", "Error"})
	defer f.Close()
	defer ef.Flush()
	for e := range errCh {
		ef.Write([]string{e.ts, e.titleID, e.pageID, e.msg})
		log.Printf("Error: %s %s", e.titleID, e.msg)
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
	of.Write([]string{
		"TimeStamp", "ID", "PageID", "Verbatim", "NameString", "OffsetStart",
		"OffsetEnd", "Odds", "Kind", "EndsNextPage",
	})
	tf.Write([]string{"ID", "Path", "PagesNumber", "NamesOccurences"})

	defer f.Close()
	defer titles.Close()
	defer of.Flush()
	defer tf.Flush()

	for t := range outCh {
		tf.Write([]string{
			t.id, t.path, strconv.Itoa(len(t.pages)), strconv.Itoa(len(t.res.Names)),
		})
		count++
		if len(t.res.Names) == 0 {
			continue
		}
		if hti.ProgressNum > 0 && count%hti.ProgressNum == 0 {
			rate := float64(count) / (time.Since(ts).Minutes())
			log.Printf("Processing %dth title. Rate %0.2f titles/min\n", count, rate)
		}
		names := generateNamesOutput(t)
		for _, n := range names {
			out := []string{
				n.timestamp, t.id, n.pageID, n.verbatim, n.nameString,
				strconv.Itoa(n.offsetStart), strconv.Itoa(n.offsetEnd),
				strconv.Itoa(int(n.odds)), n.kind, strconv.Itoa(n.endsNextPage),
			}
			of.Write(out)
		}
		if err := of.Error(); err != nil {
			log.Fatal(err)
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
func newDetectedName(p *tpage, n output.Name) *detectedName {
	var endsNextPage int
	var end int
	start := n.OffsetStart - p.offset
	if n.OffsetEnd < p.offsetNext {
		end = n.OffsetEnd - p.offset
	} else {
		end = n.OffsetEnd - p.offsetNext
		endsNextPage = 1
	}
	dn := detectedName{
		pageID:       p.id,
		verbatim:     n.Verbatim,
		nameString:   n.Name,
		offsetStart:  start,
		offsetEnd:    end,
		endsNextPage: endsNextPage,
		odds:         n.Odds,
		kind:         n.Type,
		timestamp:    ts(),
	}
	return &dn
}

// generateNamesOutput splits results by pages, instead of by title.
func generateNamesOutput(t *title) []*detectedName {
	ns := make([]*detectedName, len(t.res.Names))
	j := 0
	name := t.res.Names[j]
	for _, page := range t.pages {
		for {
			if name.OffsetStart <= page.offsetNext {
				ns[j] = newDetectedName(page, name)
				j++
				if j >= len(t.res.Names) {
					return ns
				}
				name = t.res.Names[j]
			} else {
				break
			}
		}
	}
	return ns
}
