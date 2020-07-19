package main

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	LogLevel string `json:"log_level"`
}

type RecipeList struct {
	Recipes []Recipe `json:"recipes"`
}
type Recipe struct {
	Output []Material `json:"out"`
	Input  []Material `json:"in"`
}
type Material struct {
	Item  string  `json:"i"`
	Count float32 `json:"c"`
}

type RecipeFile struct {
	Recipes []string `json:"recipes`
}

func loadJSON(file string, jsonStruct interface{}) interface{} {
	jsonFile, err := os.Open(file)
	defer jsonFile.Close()
	if err != nil {
		log.Criticalf("Could not load file %v", file)
	}
	jsonParser := json.NewDecoder(jsonFile)
	jsonParser.Decode(&jsonStruct)
	return jsonStruct
}

func loadListJSON(file string) (outRecipes RecipeList) {
	jsonFile, err := os.Open(file)
	defer jsonFile.Close()
	if err != nil {
		log.Criticalf("Could not load recipe file %v", file)
	}
	jsonParser := json.NewDecoder(jsonFile)
	tempRecipes := RecipeFile{}
	jsonParser.Decode(&tempRecipes)

	outRecipes = RecipeList{}
	for _, r := range tempRecipes.Recipes {
		mats := strings.Split(r, "<")
		outRec := mats[0]
		inRec := mats[1]
		outMats := []Material{}
		inMats := []Material{}
		for _, outMat := range strings.Split(outRec, ",") {
			outMatItem := strings.Split(outMat, "=")
			outMatItemName := outMatItem[0]
			outMatItemCount, _ := strconv.Atoi(outMatItem[1])
			outMats = append(outMats, Material{Item: outMatItemName, Count: float32(outMatItemCount)})
		}
		for _, inMat := range strings.Split(inRec, ",") {
			inMatItem := strings.Split(inMat, "=")
			inMatItemName := inMatItem[0]
			inMatItemCount, _ := strconv.Atoi(inMatItem[1])
			inMats = append(inMats, Material{Item: inMatItemName, Count: float32(inMatItemCount)})
		}
		outRecipes.Recipes = append(outRecipes.Recipes, Recipe{Output: outMats, Input: inMats})
	}
	return
}
