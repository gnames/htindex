package htindex

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"

	"github.com/gnames/gnfinder"
	"github.com/gnames/gnfinder/output"
)

type page struct {
	id         string
	offset     int
	offsetNext int
}

type title struct {
	id    string
	pages []*page
	text  []byte
	res   *output.Output
}

func (hti *HTindex) worker(inCh <-chan string, outCh chan<- *title,
	errCh chan<- error, wg *sync.WaitGroup) {
	var pages []*page
	offset := 0
	isPage, err := regexp.Compile(`\d{8}\.txt`)
	if err != nil {
		errCh <- err
	}
	defer wg.Done()
	opts := []gnfinder.Option{
		gnfinder.OptDict(hti.dict),
		gnfinder.OptBayes(true),
	}
	gnf := gnfinder.NewGNfinder(opts...)
	_ = gnf
	for file := range inCh {
		t := title{id: getID(file), pages: pages}
		r, err := zip.OpenReader(hti.rootPrefix + file)
		if err != nil {
			errCh <- err
		}
		for _, f := range r.File {
			fn := f.Name
			fnl := len(fn)
			if fnl < 12 || !isPage.MatchString(fn[fnl-12:fnl]) {
				continue
			}
			// fnEnd := fn[fnl-12 : fnl-4]
			rc, err := f.Open()
			if err != nil {
				errCh <- err
			}
			bs, err := ioutil.ReadAll(rc)
			if err != nil {
				errCh <- err
			}
			p := page{id: fn[fnl-12 : fnl-4], offset: offset}
			pageUTF := []rune(string(bs))
			offset += len(pageUTF)
			p.offsetNext = offset
			t.pages = append(t.pages, &p)
			t.text = append(t.text, bs...)
			rc.Close()
		}
		r.Close()
		t.res = gnf.FindNames(t.text)
		outCh <- &t
	}
}

func getID(p string) string {
	el := strings.Split(p, "/")
	return fmt.Sprintf("%s.%s", el[0], el[len(el)-2])
}
