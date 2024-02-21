package compare

import (
	"fmt"
)

func Compare(old, new *DBReader) {
	oldRecipe := (*old).getRecipe()
	newRecipe := (*new).getRecipe()
	compareNames(oldRecipe, newRecipe)
	compareTime(oldRecipe, newRecipe)
	compareIngredients(oldRecipe, newRecipe)
}

func compareNames(oldRecipe, newRecipe *Recipe) {
	i := 0
	for _, old := range oldRecipe.Cake {
		new := newRecipe.Cake[i]
		if old.Name != new.Name {
			fmt.Printf("Added cake \"%s\"\n", new.Name)
			fmt.Printf("REMOVED cake \"%s\"\n", old.Name)
		}
		i++
	}
}

func compareTime(oldRecipe, newRecipe *Recipe) {
	i := 0
	for _, old := range oldRecipe.Cake {
		new := newRecipe.Cake[i]
		if old.Time != new.Time {
			fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", old.Name, new.Time, old.Time)
		}
		i++
	}
}

func compareIngredients(oldRecipe, newRecipe *Recipe) {
	for _, newRecipe := range newRecipe.Cake {
		for _, oldRecipe := range oldRecipe.Cake {
			if newRecipe.Name == oldRecipe.Name {
				status := "ADDED"
				compareNamesIngredients(&oldRecipe.CakeIngredients, &newRecipe.CakeIngredients, newRecipe.Name, status)
				status = "REMOVED"
				compareNamesIngredients(&newRecipe.CakeIngredients, &oldRecipe.CakeIngredients, newRecipe.Name, status)
				compareAllIngredients(&oldRecipe.CakeIngredients, &newRecipe.CakeIngredients, newRecipe.Name)
				break
			}
		}
	}
}

func compareNamesIngredients(old, new *[]Ingredients, cakeName, status string) {
	for _, ingredientNew := range *new {
		found := false
		for _, ingredientOld := range *old {
			if ingredientOld.Name == ingredientNew.Name {
				found = true
			}
		}
		if !found {
			fmt.Printf("%s ingredient \"%s\" for cake  \"%s\"\n", status, ingredientNew.Name, cakeName)
		}
	}
}

func compareAllIngredients(old, new *[]Ingredients, cakeName string) {
	for _, ingredientNew := range *new {
		for _, ingredientOld := range *old {
			if ingredientOld.Name == ingredientNew.Name {
				oldName := ingredientOld.Name
				if len(ingredientNew.Unit) > 0 && ingredientOld.Unit != ingredientNew.Unit {
					fmt.Printf("CHANGED unit for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", oldName, cakeName, ingredientNew.Unit, ingredientOld.Unit)
				} else if ingredientOld.Count != ingredientNew.Count {
					fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", oldName, cakeName, ingredientNew.Count, ingredientOld.Count)
				} else if len(ingredientNew.Unit) == 0 && ingredientOld.Unit != ingredientNew.Unit {
					fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake  \"%s\"\n", ingredientOld.Unit, oldName, cakeName)
				}
			}
		}
	}
}
