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
var isPage = regexp.MustCompile(`\d{8}\.txt`)

// pageContent allows to presort pages from zip file in case if they are given
// in a wrong order.
type pageContent struct {
	id   string
	text []byte
}

// tpage represents metadata of a page from a title.
type tpage struct {
	id         string
	offset     int
	offsetNext int
}

// title represents data and metadata from a title/book/volume.
type title struct {
	id    string
	path  string
	pages []*tpage
	text  []byte
	res   *output.Output
}

// byID allows to sort pageContent slice using its `id` field.
type byID []*pageContent

func (b byID) Len() int           { return len(b) }
func (b byID) Less(i, j int) bool { return b[i].id < b[j].id }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

// worker is the maing workhorse of the app. It reads zip file, extracts
// data from pages, and prepares title data and results of name-finding for
// the output. In case if some errors happened during processing, they will be
// prepared for logging.
func (hti *HTindex) worker(inCh <-chan string, outCh chan<- *title,
	errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	opts := []gnfinder.Option{
		gnfinder.OptDict(hti.dict),
		gnfinder.OptBayes(true),
	}
	gnf := gnfinder.NewGNfinder(opts...)

	for zipPath := range inCh {
		offset := 0
		var pages []*tpage
		t := title{id: getID(zipPath), path: zipPath, pages: pages}
		r, err := zip.OpenReader(filepath.Join(hti.rootPrefix, zipPath))
		if err != nil {
			errCh <- err
		}

		for _, pc := range pagesContent(r, errCh) {
			p := tpage{id: pc.id, offset: offset}
			pageUTF := []rune(string(pc.text))
			offset += len(pageUTF)
			p.offsetNext = offset
			t.pages = append(t.pages, &p)
			t.text = append(t.text, pc.text...)
		}
		r.Close()
		t.res = gnf.FindNames(t.text)
		outCh <- &t
	}
}

// pagesContent generates a list of all pages with their texts sorted according
// to their position in the title.
func pagesContent(r *zip.ReadCloser, errCh chan<- error) []*pageContent {
	var pages []*pageContent
	for _, f := range r.File {
		fn := f.Name
		fnl := len(fn)
		if fnl < 12 || !isPage.MatchString(fn[fnl-12:fnl]) {
			continue
		}
		zf, err := f.Open()
		if err != nil {
			errCh <- err
		}
		id := fn[fnl-12 : fnl-4]
		text, err := ioutil.ReadAll(zf)
		if err != nil {
			errCh <- err
		}
		pages = append(pages, &pageContent{id: id, text: text})
		zf.Close()
	}
	sort.Sort(byID(pages))
	return pages
}

// getID generates the id of a title from its filepath.
func getID(p string) string {
	el := strings.Split(p, "/")
	return fmt.Sprintf("%s.%s", el[0], el[len(el)-2])
}
