package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	err := checkArgs()
	if err != nil {
		log.Fatal(err)
	}
	oldFile, err := openFile(os.Args[2])

	if err != nil {
		log.Fatal(err)
	}
	defer oldFile.Close()

	newFile, err := openFile(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()
	old := fileToMap(oldFile)
	new := fileToMap(newFile)
	compareFiles(old, new)

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

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func fileToMap(file *os.File) map[string]string {
	strs := make(map[string]string)
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		line := reader.Text()
		strs[line] = line
	}
	return strs
}

func compareFiles(old, new map[string]string) {
	for key, value := range new {
		if _, ok := old[key]; !ok {
			fmt.Println("ADDED:", value)
		}
	}

	for key, value := range old {
		if _, ok := new[key]; !ok {
			fmt.Println("REMOVED:", value)
		}
	}
}
