package compare

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

type Recipe struct {
	Cake []struct {
		Name            string        `json:"name" xml:"name"`
		Time            string        `json:"time" xml:"stovetime"`
		CakeIngredients []Ingredients `json:"ingredients" xml:"ingredients>item"`
	} `json:"cake" xml:"cake"`
}

type Ingredients struct {
	Name  string `json:"ingredient_name" xml:"itemname"`
	Count string `json:"ingredient_count" xml:"itemcount"`
	Unit  string `json:"ingredient_unit,omitempty" xml:"itemunit,omitempty"`
}

type DBReader interface {
	ReadRecipe() error
	PrintRecipe()
	getRecipe() *Recipe
}

type StoreJson struct {
	Data   []byte
	Recipe *Recipe
}

type StoreXml struct {
	Data   []byte
	Recipe *Recipe
}

func (j *StoreJson) ReadRecipe() error {
	j.Recipe = &Recipe{}
	err := json.Unmarshal(j.Data, j.Recipe)
	if err != nil {
		return err
	}
	return nil
}

func (j *StoreJson) PrintRecipe() {
	res, err := xml.MarshalIndent(j.Recipe, "", "    ")
	if err == nil {
		fmt.Println(string(res))
	}

}

func (x *StoreXml) ReadRecipe() error {
	x.Recipe = &Recipe{}
	err := xml.Unmarshal(x.Data, x.Recipe)
	if err != nil {
		return err
	}
	return nil
}

func (x *StoreXml) PrintRecipe() {
	res, err := json.MarshalIndent(x.Recipe, "", "    ")
	if err == nil {
		fmt.Println(string(res))
	}
}

func (x *StoreXml) getRecipe() *Recipe {
	return x.Recipe
}

func (j *StoreJson) getRecipe() *Recipe {
	return j.Recipe
}
