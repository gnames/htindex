package htindex

import (
	"archive/zip"
	"bufio"
	"fmt"
	"os"
)

func (hti *HTindex) Run() error {
	inCh := make(chan string)
	errCh := make(chan error)
	go outputError(errCh)
	go hti.outputInput(inCh, errCh)
	if err := hti.readInput(inCh, errCh); err != nil {
		return err
	}
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

func outputError(errCh <-chan error) {
	for err := range errCh {
		fmt.Println(err)
	}
}

func (hti *HTindex) outputInput(inCh <-chan string, errCh chan<- error) {
	for file := range inCh {
		fmt.Println(file)
		fmt.Println()
		r, err := zip.OpenReader(hti.rootPrefix + file)
		if err != nil {
			errCh <- err
		}
		defer r.Close()
		for _, f := range r.File {
			fmt.Printf("File: %s\n", f.Name)
			rc, err := f.Open()
			if err != nil {
				errCh <- err
			}

			scanner := bufio.NewScanner(rc)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				errCh <- err
			}
			rc.Close()
		}
	}
}
