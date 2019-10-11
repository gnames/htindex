package htindex

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/output"
)

// isPage determines if a file represents a page with text from the title.
var isPage = regexp.MustCompile(`\d{6}\.txt`)

// page allows to presort pages from zip file in case if they are given
// in a wrong order.
type page struct {
	id   string
	text []byte
	res  *output.Output
}

// title represents data and metadata from a title/book/volume.
type title struct {
	id               string
	path             string
	pages            []*page
	namesNum         int
	pagesNumBadNames int
}

type htiError struct {
	ts      string
	titleID string
	pageID  string
	msg     string
}

// byID allows to sort page slice using its `id` field.
type byID []page

func (b byID) Len() int           { return len(b) }
func (b byID) Less(i, j int) bool { return b[i].id < b[j].id }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// worker is the maing workhorse of the app. It reads zip file, extracts
// data from pages, and prepares title data and results of name-finding for
// the output. In case if some errors happened during processing, they will be
// prepared for logging.
func (hti *HTindex) worker(inCh <-chan string, outCh chan<- *title,
	errCh chan<- *htiError, wg *sync.WaitGroup) {
	defer wg.Done()

	opts := []gnfinder.Option{
		gnfinder.OptDict(hti.Dict),
		gnfinder.OptBayes(true),
	}
	gnf := gnfinder.NewGNfinder(opts...)

	for zipPath := range inCh {
		var pages []*page
		t := title{id: getID(zipPath), path: zipPath, pages: pages}
		r, err := zip.OpenReader(filepath.Join(hti.RootPrefix, zipPath))
		if err != nil {
			errCh <- &htiError{msg: err.Error(), titleID: t.id, ts: ts()}
		}
		pcs, pagesNumBadNames := pagesContent(r, errCh)
		if pagesNumBadNames > 0 {
			msg := fmt.Sprintf("non-standard naming for %d pages", pagesNumBadNames)
			errCh <- &htiError{msg: msg, titleID: t.id, ts: ts()}
		}
		t.pagesNumBadNames = pagesNumBadNames
		if len(pcs) == 0 {
			errCh <- &htiError{ts: ts(), titleID: t.id, msg: "no pages detected"}
			continue
		}

		for _, pc := range pcs {
			pc.res = gnf.FindNames(pc.text)
			t.namesNum += len(pc.res.Names)
			t.pages = append(t.pages, &pc)
		}
		r.Close()
		outCh <- &t
	}
}

// pagesContent generates a list of all pages with their texts sorted according
// to their position in the title.
func pagesContent(r *zip.ReadCloser, errCh chan<- *htiError) ([]page, int) {
	badPageName := 0
	var pages []page
	for _, f := range r.File {
		fn := f.Name
		fnl := len(fn)
		if fnl < 12 || !isPage.MatchString(fn[fnl-12:fnl]) {
			continue
		}
		zf, err := f.Open()
		if err != nil {
			errCh <- &htiError{msg: err.Error()}
		}
		id := fn[fnl-12 : fnl-4]
		if !strings.HasPrefix(id, "00") {
			badPageName++
		}
		text, err := ioutil.ReadAll(zf)
		if err != nil {
			errCh <- &htiError{msg: err.Error()}
		}
		pages = append(pages, page{id: id, text: text})
		zf.Close()
	}
	sort.Sort(byID(pages))
	return pages, badPageName
}

// getID generates the id of a title from its filepath.
func getID(p string) string {
	el := strings.Split(p, "/")
	return fmt.Sprintf("%s.%s", el[0], el[len(el)-2])
}
