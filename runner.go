package htindex

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Run is the main method for creation of the scientific names index.
func (hti *HTindex) Run() error {
	fmt.Printf("Processing with %d 'threads'\n", hti.JobsNum)
	inCh := make(chan string)
	errCh := make(chan *htiError)
	outCh := make(chan *title)
	var wg sync.WaitGroup
	var wgOut sync.WaitGroup
	wg.Add(hti.JobsNum)
	wgOut.Add(2)
	go hti.outputError(errCh, &wgOut)
	go hti.outputResult(outCh, &wgOut)
	for i := 0; i < hti.JobsNum; i++ {
		go hti.worker(inCh, outCh, errCh, &wg)
	}
	if err := hti.readInput(inCh, errCh); err != nil {
		return err
	}
	wg.Wait()
	close(outCh)
	close(errCh)
	wgOut.Wait()
	return nil
}

// readInput traverses the input file and sends paths to title's zip files to
// further processes.
func (hti *HTindex) readInput(inCh chan<- string,
	errCh chan<- *htiError) error {
	file, err := os.Open(hti.InputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	defer close(inCh)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inCh <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		errCh <- &htiError{msg: err.Error()}
	}
	return nil
}
