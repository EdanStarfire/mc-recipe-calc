# mc-recipe-calc
Calculator for minecraft recipes

**Don't forget to populate a recipes file with your recipes or this will not work!**

## Usage
```bash
mc-recipe-calc recipes.json targetitem targetcount
    recipelist.json      - File used to load the recipe list
    targetitem           - Item you want to make X of
    targetcount          - How many items you want to make
```
Example:
```bash
mc-recipe-calc manufactio.json elite_circuit 64
```
Output:
```bash
[64x] elite_circuit -> needs: [{copper_ingot 2432} {iron_ingot 1728} {redstone 1536} {plastic_sheet 256} {gold_wire 64} {sulfur 960}]
```

## Recipe File
Recipe file is a json file of the format:

```json
{
  "recipes": [
    "outputitem=1<inputitem1=4,inputitem2=2"
    ,"outputitem2=6,outputitem=1<output=2,intputitem1=4,inputitem3=1"
  ]
}
```
*The file `manufactio.json` is provided as an example of how to input recipes.*

In this case the formatting is item=count (where count is the number of items being made or needed).

It can support multiple items on the input or output.

When it cannot find a recipe for a specific input item, it just assumes that this is a raw material. A good example of this would be if you already have pistons being made automatically, and you don't want pistons to be broken down into their inputs, you can just not have a recipe that outputs pistons in your recipe list and it'll result in the pistons being treated as an input item.
