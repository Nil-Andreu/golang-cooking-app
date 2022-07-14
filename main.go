package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Define Recipe Data Structure
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
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

// Define the Routes Handlers
func GetSliceRecipes(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// Definintion for the creation of a new recipe
func NewRecipeHandler(c *gin.Context) {
	// Define a new variable which is recipe, based on Recipe object
	var recipe Recipe

	// The ShouldBindJSON convert the request body into Recipe struct and assigns unique identifier with an external package called xid
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

// To update the Recipe
func UpdateRecipeHandler(c *gin.Context) {
	// Obtain first which is the parameter
	id := c.Param("id")

	// Initialize the recipe
	var recipe Recipe

	// Bind the recipe from json
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error" : err.Error()}) 
		return
	}

	// We obtain the recipe we looked for
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusBadRequest, gin.H {
			"error" : "Index not found",
		})
	}
}

func main() {
	router := gin.Default()
	router.GET("/recipes", GetSliceRecipes)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.Run(":6000")
}
