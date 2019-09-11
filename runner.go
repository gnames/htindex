package htindex

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func (hti *HTindex) Run() error {
	fmt.Println(hti)
	inCh := make(chan string)
	errCh := make(chan error)
	outCh := make(chan *title)
	var wg sync.WaitGroup
	var wgOut sync.WaitGroup
	wg.Add(hti.jobsNum)
	wgOut.Add(1)
	go hti.outputError(errCh)
	go hti.outputResult(outCh, &wgOut)
	for i := 0; i < hti.jobsNum; i++ {
		go hti.worker(inCh, outCh, errCh, &wg)
	}
	if err := hti.readInput(inCh, errCh); err != nil {
		return err
	}
	wg.Wait()
	close(outCh)
	wgOut.Wait()
	return nil
}

func (hti *HTindex) readInput(inCh chan<- string, errCh chan<- error) error {
	file, err := os.Open(hti.inputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	defer close(inCh)
	defer close(errCh)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		inCh <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		errCh <- err
	}
	return nil
}
