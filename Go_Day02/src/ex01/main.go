package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"unicode/utf8"
)

type Flags struct {
	l, w, m bool
}

var fl Flags

func init() {
	flag.BoolVar(&fl.l, "l", false, "lines")
	flag.BoolVar(&fl.w, "w", false, "words")
	flag.BoolVar(&fl.m, "m", false, "characters")
}

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	err := checkArgs()
	if err != nil {
		log.Fatal(err)
	}
	files := flag.Args()
	for _, file := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			err = fileProcessing(file)
			if err != nil {
				log.Fatal(err)
			}
		}(file)
	}
	wg.Wait()
}

func fileProcessing(file string) error {
	f, err := openFile(file)
	if err != nil {
		return err
	}
	if fl.l {
		amountLines, err := counting(f, '\n')
		if err != nil {
			return err
		}
		fmt.Printf("%d\t%s\n", amountLines, file)
	} else if fl.m {
		fileScanner := bufio.NewScanner(f)
		amountChars := 0
		for fileScanner.Scan() {
			str := fileScanner.Text()
			amountChars += utf8.RuneCountInString(str) + 1
		}
		fmt.Printf("%d\t%s\n", amountChars-1, file)
	} else if fl.w {
		amountWords, err := counting(f, ' ')
		if err != nil {
			return err
		}
		fmt.Printf("%d\t%s\n", amountWords+1, file)
	}

	return nil
}

func checkArgs() error {
	amount := flag.NFlag()
	if amount > 1 {
		return fmt.Errorf("the program expect a one flag per input")
	}

	if amount == 0 {
		fl.w = true
	}

	if amountFiles := flag.NArg(); amountFiles < 1 {
		return fmt.Errorf("not enough files")
	}
	return nil
}

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func counting(r io.Reader, s byte) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	sep := []byte{s}
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], sep)
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
