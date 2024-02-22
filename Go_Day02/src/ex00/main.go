package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var bufout = bufio.NewWriter(os.Stdout)

type Flags struct {
	file, dir, sl bool
	extension     string
}

var fl Flags

func main() {
	initFlags(&fl)
	flag.Parse()
	paths, err := checkArguments()

	if err != nil {
		log.Fatal(err.Error())
	}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		err = filepath.Walk(path, searchFiles)
		if err != nil {
			log.Fatal(err)
		}
	}
	bufout.Flush()
}

func checkArguments() ([]string, error) {
	//Функция возвращает слайс аргументов НЕ флагов
	args := flag.Args()

	if len(args) < 1 {
		return nil, errors.New("missing file path")
	}

	if !fl.file && fl.extension != "" {
		return nil, errors.New("the -ext flag must be used with the -f flag")
	}

	if !fl.file && !fl.dir && !fl.sl {
		fl.file, fl.dir, fl.sl = true, true, true
	}
	return args, nil
}

func initFlags(fl *Flags) {
	flag.BoolVar(&fl.file, "f", false, "file")
	flag.BoolVar(&fl.dir, "d", false, "directory")
	flag.BoolVar(&fl.sl, "sl", false, "symbolic links")
	flag.StringVar(&fl.extension, "ext", "", "extension")
}

func searchFiles(wPath string, info os.FileInfo, err error) error {
	if info.Mode().IsRegular() && fl.extension != "" {

		filextension := filepath.Ext(wPath)

		if filextension == ("." + fl.extension) {
			fmt.Fprintln(bufout, wPath)
		}
	} else {

		if fl.sl {
			if link, _ := filepath.EvalSymlinks(wPath); link != wPath {
				if _, err := os.Stat(wPath); err == nil {
					fmt.Fprintln(bufout, wPath, "->", link)
				} else {
					fmt.Fprintln(bufout, wPath, "->", "[broken]")
				}
			}
		}

		if fl.file && info.Mode().IsRegular() {
			fmt.Fprintln(bufout, wPath)
		}

		if fl.dir && info.IsDir() {
			fmt.Fprintln(bufout, wPath)
		}
	}

	return nil
}
