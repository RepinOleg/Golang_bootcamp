package main

import (
	"flag"
	"fmt"
	"log"
	"myRotate/archiever"
	"os"
	"sync"
)

var path *string

func init() {
	path = flag.String("a", ".", "path")
	flag.Parse()
	err := checkArgs()
	if err != nil {
		log.Fatal(err)
	}
	err = checkPath(*path)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	files := flag.Args()
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go archiever.CreateArchive(*path, file, &wg)
	}
	wg.Wait()
}

func checkArgs() error {
	if flag.NArg() < 1 {
		return fmt.Errorf("not enough log files")
	}

	return nil
}

func checkPath(path string) error {
	if path != "." {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return fmt.Errorf("%s is not a directory", path)
		}
	}
	return nil
}
