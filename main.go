package main

import (
	"time"
	"net/http"
	"github.com/rs/xid"
	"github.com/gin-gonic/gin"
)

// Define Recipe Data Structure
type Recipe struct {
	ID				string		`json:"id"`	
	Name 			string 		`json:"name"`
	Tags			[]string 	`json:"ingredients"`
	Instructions	[]string	`json:"instructions"`
	PublishedAt		time.Time	`json:"publishedAt"`
}

// Define recipes new list (later we will put it in a data base)
var recipes []Recipe
func init() {						// Init function is called when the app is initialized
	recipes = make([]Recipe, 0) 	// Create the Recipe, with length of 0
}

// Define the Routes Handlers
func GetRecipes(c *gin.Context) {
	c.JSON(200, gin.H {

	})
}

// Definintion for the creation of a new recipe
func NewRecipeHandler(c *gin.Context) {
	// Define a new variable which is recipe, based on Recipe object
	var recipe Recipe

	// The ShouldBindJSON convert the request body into Recipe struct and assigns unique identifier with an external package called xid
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"error": err.Error()}) 
		return 
		// Function does not return anything
	}

	// Create a new identifier with xid
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe) // Append to the list of recipes, the new recipe

	c.JSON(http.StatusOK, recipe)
}



func main() {
	router := gin.Default()
	router.GET("/", GetRecipes)
	router.POST("/recipes", NewRecipeHandler)
	router.Run()
}

