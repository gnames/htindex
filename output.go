package htindex

import (
	"fmt"
	"sync"
)

func (hti *HTindex) outputError(errCh <-chan error) {
	for err := range errCh {
		fmt.Println(err)
	}
}

func (hti *HTindex) outputResult(outCh <-chan *title, wgOut *sync.WaitGroup) {
	defer wgOut.Done()
	for t := range outCh {
		fmt.Printf("%s pages:%d names:%d\n", t.id, len(t.pages), len(t.res.Names))
	}
}
