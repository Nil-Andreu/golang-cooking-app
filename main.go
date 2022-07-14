// Recipes
// A simple API for adding, removing, reviewing and updating recipes.
// Schemes: http
// Host: localhost:6000
// BasePath: /
// Version: 1.0.0
// Contact: Nil Andreu <nilandreug@gmail.com>
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Define Recipe Data Structure
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tag          []string  `json:"tag"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

// Define recipes new list (later we will put it in a data base)
var recipes []Recipe

func init() { // Init function is called when the app is initialized
	recipes = make([]Recipe, 0) // Create the Recipe, with length of 0
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//  '200':
//	  description: Successful operation
func GetSliceRecipes(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// swagger:operation POST /recipes recipe newRecipe
// Create a new recipe
// ---
// parameters:
// - name: name
//   in: body
//   description: The name of the Recipe
//   require: true
//   type: string
// - name: tag
//   in: body
//   description: Tags for the recipe
//   require: true
//   schema: 
//    type: array
//    items:
//     type: string
//    example: ["tag1", "tag2"]
// - name: ingredients
//   in: body
//   description: The ingredients of the Recipe
//   require: true
//   schema: 
//    type: array
//    items:
//     type: string
//    example: ["ingredient1", "ingredient2"]
// - name: instructions
//   in: body
//   description: The instructions of the Recipe
//   require: true
//   type: list
//   schema: 
//    type: array
//    items:
//     type: string
//    example: ["instruction1", "instruction2"]
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//  '200':
//	  description: Successful operation
//  '404':
//    description: Invalid Recipe ID
func NewRecipeHandler(c *gin.Context) {
	// Define a new variable which is recipe, based on Recipe object
	var recipe Recipe

	// The ShouldBindJSON convert the request body into Recipe struct in the recipe variable
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
		// Function does not return anything
	}

	// Create a new identifier with xid
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe) // Append to the list of recipes, the new recipe

	// We write the new information
	file, _ := json.MarshalIndent(recipes, "", "")   // Write recipes as json structure in bytes
	_ = ioutil.WriteFile("recipes.json", file, 0644) // store it

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /recipes/:id recipe updateRecipe
// To update the Recipe
// ---
// produces:
// - application/json
// consumes:
// - application/json
// parameters:
// - name: id
//   in: path
//   description: The id of the Recipe we want to update
//   require: true
//   type: int
// - name: ID
//   in: body
//   description: The new ID for the Recipe (optional)
//   require: false
//   type: string
// - name: name
//   in: body
//   description: The name of the Recipe
//   require: true
//   type: string
// - name: tag
//   in: body
//   description: Tags for the recipe
//   require: true
//   schema: 
//    type: array
//    items:
//     type: string
//    example: ["tag1", "tag2"]
// - name: ingredients
//   in: body
//   description: The ingredients of the Recipe
//   require: true
//   schema: 
//    type: array
//    items:
//     type: string
//    example: ["ingredient1", "ingredient2"]
// - name: instructions
//   in: body
//   description: The instructions of the Recipe
//   require: true
//   schema: 
//    type: array
//    items:
//     type: string
//    example: ["instruction1", "instruction2"]
// responses:
//  '200':
//	  description: Successful operation
//  '400':
//    description: Invalid Input
//  '404':
//    description: Invalid id
func UpdateRecipeHandler(c *gin.Context) {
	// Obtain first which is the parameter
	id := c.Param("id")

	// Initialize the recipe
	var recipe Recipe

	// Bind the recipe of the body of the request that is formatted in json to &recipe variable
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	// In the case that we have not provided with an id
	if recipe.ID == "" {
		recipe.ID = xid.New().String()
	}

	// We obtain the recipe we looked for
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	// In the case we did not found it
	if index == -1 {
		message := fmt.Sprintf("Index %s not found", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
	}

	// If we found it, replace the recipe at that position
	recipes[index] = recipe

	// And we now need to re store it on the json file
	file, _ := json.MarshalIndent(recipes, "", "") // Write recipes as json structure in bytes
	_ = ioutil.WriteFile("recipes.json", file, 0644)

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation DELETE /recipes/:id recipes deleteRecipes
// Delete a given recipe
// ---
// parameters:
// - name: id
//   in: path
//   description: The recipe id we want to remove
//   require: true
//   type: string 
// produces:
// - application/json
// responses:
//  '200':
//   description: Successful operation
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	// Initialize the recipe
	var recipe Recipe

	// Bind the recipe of the body of the request that is formatted in json to &recipe variable
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		message := fmt.Sprintf("Index %s not found", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": message,
		})
	}

	// And now we delete the recipe at that position
	recipes = append(recipes[:index], recipes[index+1:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Deleted correctly",
	})

}

// swagger:operation GET /recipes/search recipes searchRecipes
// Search recipes based on tag
// ---
// parameters:
// - name: tag
//   in: path
//   required: true
//   type: string
// produces:
// - application/json
// responses:
//  '200':
//   description: Successful operation
func SearchRecipeHandler(c *gin.Context) {
	Tag := c.Query("tag")
	resultRecipes := make([]Recipe, 0)

	for recipeIndex := 0; recipeIndex < len(recipes); recipeIndex++ {
		// And for this recipe, we will look for all the tags that it has
		recipe := recipes[recipeIndex]

		for tagIndex := 0; tagIndex < len(recipe.Tag); tagIndex++ {
			// In the case that those are equal
			if strings.EqualFold(recipe.Tag[tagIndex], Tag) {
				resultRecipes = append(resultRecipes, recipe)
			}
		}

		// In the case there is no match, handle this situation
		if len(resultRecipes) == 0 {
			message := fmt.Sprintf("Do not found any recipe with the tag %s", Tag)
			c.JSON(http.StatusOK, gin.H{
				"message": message,
			})
		} else {
			c.JSON(http.StatusOK, resultRecipes)
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/recipes", GetSliceRecipes)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run(":6000")
}
