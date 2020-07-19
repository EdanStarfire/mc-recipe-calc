package main

import (
	"os"
	"strconv"

	"github.com/juju/loggo"
	"github.com/pkg/errors"
)

var log = loggo.GetLogger("")

type options struct {
	RecipeFile  string
	TargetItem  string
	TargetCount float32
}

func main() {
	cfg := loadJSON("config.json", &Config{}).(*Config)
	loggo.ConfigureLoggers(cfg.LogLevel)
	if log.LogLevel() != loggo.DEBUG && log.LogLevel() != loggo.TRACE {
		loggo.ReplaceDefaultWriter(loggo.NewSimpleWriter(os.Stdout, infoLoggoWriter))
	}
	opts := parseArgs()

	//recipeList := loadJSON("manufactio.json", &RecipeList{}).(*RecipeList)
	recipeList := loadListJSON(opts.RecipeFile)

	for _, recipe := range recipeList.Recipes {
		log.Tracef("Out -> %v ", recipe.Output)
	}
	count := opts.TargetCount
	target := opts.TargetItem
	inMats := getRawMaterials(target, count, &recipeList)
	log.Infof("[%vx] %v -> needs: %v", count, target, inMats)
}

func parseArgs() (opts options) {
	opts = options{}
	args := os.Args[1:]
	if len(args) != 3 {
		log.Errorf("invalid runtime options")
		printUsage()
		os.Exit(1)
	}
	opts.RecipeFile = args[0]
	opts.TargetItem = args[1]
	targetCount, err := strconv.ParseFloat(args[2], 32)
	if err != nil {
		log.Errorf("couldn't parse %v into a number", args[2])
		printUsage()
		os.Exit(1)
	}
	opts.TargetCount = float32(targetCount)
	return
}

func printUsage() {
	log.Infof("")
	log.Infof("Usage:")
	log.Infof("    mc-recipe-calc recipes.json targetitem targetcount")
	log.Infof("        recipes.json         - File used to load the recipe list")
	log.Infof("        targetitem           - Item you want to make X of")
	log.Infof("        targetcount          - How many items you want to make")
	log.Infof("")
}

func getRecipe(target string, list *RecipeList) (rec Recipe, err error) {
	rec = Recipe{}
	for _, recipe := range list.Recipes {
		for _, outMaterial := range recipe.Output {
			if outMaterial.Item == target {
				rec = recipe
				return
			}
		}
	}
	err = errors.New("no recipe found")
	return
}

func getRawMaterials(target string, count float32, list *RecipeList) (neededMats []Material) {
	rec, err := getRecipe(target, list)
	if err != nil {
		log.Tracef("raw material: [%vx] %v", count, target)
		neededMats = appendMats(neededMats, []Material{{Item: target, Count: float32(count)}})
		return
	}
	multiplier := count
	for _, outMat := range rec.Output {
		if outMat.Item == target {
			multiplier = count / outMat.Count
		}
	}
	for _, inMat := range rec.Input {
		targetCount := inMat.Count * multiplier
		inMats := getRawMaterials(inMat.Item, targetCount, list)
		neededMats = appendMats(neededMats, inMats)
	}
	return
}

func appendMats(currentMats []Material, newMats []Material) (outMats []Material) {
	outMats = []Material{}
	log.Tracef("appending mats %+v  +++  %v", currentMats, newMats)
	for _, cMat := range currentMats {
		matCount := cMat.Count
		for _, nMat := range newMats {
			if nMat.Item == cMat.Item {
				matCount += nMat.Count
			}
		}
		outMats = append(outMats, Material{cMat.Item, matCount})
	}
	for _, nMat := range newMats {
		foundMat := false
		for _, cMat := range currentMats {
			if nMat.Item == cMat.Item {
				foundMat = true
				break
			}
		}
		if !foundMat {
			outMats = append(outMats, nMat)
		}
	}
	return
}

func infoLoggoWriter(entry loggo.Entry) string {
	return entry.Message
}
