package main

import (
	"compareDB/compare"
	"fmt"
	"os"
	"strings"
)

func main() {
	err := checkArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	filenameOldRecipe := os.Args[2]
	filenameNewRecipe := os.Args[len(os.Args)-1]

	arrData1, err := readFile(filenameOldRecipe)
	arrData2, err2 := readFile(filenameNewRecipe)
	if err != nil {
		fmt.Println(err)
		return
	} else if err2 != nil {
		fmt.Println(err2)
		return
	}

	readerOld, err := createStores(filenameOldRecipe, arrData1)
	readerNew, err2 := createStores(filenameNewRecipe, arrData2)

	if err != nil {
		fmt.Println(err)
		return
	} else if err2 != nil {
		fmt.Println(err2)
		return
	}

	err = readerOld.ReadRecipe()
	err2 = readerNew.ReadRecipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	compare.Compare(&readerOld, &readerNew)
}

func checkArgs() error {
	if len(os.Args) != 5 {
		return fmt.Errorf("the number of arguments must be five")
	}

	if os.Args[1] != "--old" {
		return fmt.Errorf("the second parameter should be called --old")
	}

	if os.Args[3] != "--new" {
		return fmt.Errorf("the fourth parameter should be called --new")
	}

	return nil
}

func readFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func createStores(filename string, data []byte) (compare.DBReader, error) {
	var reader compare.DBReader
	if strings.HasSuffix(filename, "xml") {
		reader = &compare.StoreXml{Data: data}
	} else if strings.HasSuffix(filename, "json") {
		reader = &compare.StoreJson{Data: data}
	} else {
		return nil, fmt.Errorf("wrong extension")
	}
	return reader, nil
}
