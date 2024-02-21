package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Recipe struct {
	Cake []struct {
		Name        string `json:"name" xml:"name"`
		Time        string `json:"time" xml:"stovetime"`
		Ingredients []struct {
			Name  string `json:"ingredient_name" xml:"itemname"`
			Count string `json:"ingredient_count" xml:"itemcount"`
			Unit  string `json:"ingredient_unit,omitempty" xml:"itemunit,omitempty"`
		} `json:"ingredients" xml:"ingredients>item"`
	} `json:"cake" xml:"cake"`
}

type DBReader interface {
	ReadRecipe() error
	PrintRecipe()
}

type storeJson struct {
	file   []byte
	Recipe *Recipe
}

type storeXml struct {
	file   []byte
	Recipe *Recipe
}

func (j *storeJson) ReadRecipe() error {
	j.Recipe = &Recipe{}
	err := json.Unmarshal(j.file, j.Recipe)
	if err != nil {
		return err
	}
	return nil
}

func (j *storeJson) PrintRecipe() {
	res, err := xml.MarshalIndent(j.Recipe, "", "    ")
	if err == nil {
		fmt.Println(string(res))
	}

}

func (x *storeXml) ReadRecipe() error {
	x.Recipe = &Recipe{}
	err := xml.Unmarshal(x.file, x.Recipe)
	if err != nil {
		return err
	}
	return nil
}

func (j *storeXml) PrintRecipe() {
	res, err := json.MarshalIndent(j.Recipe, "", "    ")
	if err == nil {
		fmt.Println(string(res))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}

	var reader DBReader
	filename := os.Args[len(os.Args)-1]
	file, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.HasSuffix(filename, "xml") {
		reader = &storeXml{file: file}
	} else if strings.HasSuffix(filename, "json") {
		reader = &storeJson{file: file}
	} else {
		fmt.Println("Wrong extension")
		return
	}

	err = reader.ReadRecipe()

	if err != nil {
		fmt.Println(err)
		return
	}

	reader.PrintRecipe()

}
