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

var isPage = regexp.MustCompile(`\d{8}\.txt`)

type pageContent struct {
	id   string
	text []byte
}

type tpage struct {
	id         string
	offset     int
	offsetNext int
}

type title struct {
	id    string
	pages []*tpage
	text  []byte
	res   *output.Output
}

type byID []*pageContent

func (b byID) Len() int           { return len(b) }
func (b byID) Less(i, j int) bool { return b[i].id < b[j].id }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

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
		t := title{id: getID(zipPath), pages: pages}
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

func getID(p string) string {
	el := strings.Split(p, "/")
	return fmt.Sprintf("%s.%s", el[0], el[len(el)-2])
}
