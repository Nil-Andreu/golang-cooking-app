# Cooking App
Will create a cooking application with Golang.

In this application, you will be able to handle recipes.

The object of the recipes is the following:
```
    type Recipe struct {
        ID				string		`json:"id"`	
        Name 			string 		`json:"name"`
        Tags			[]string 	`json:"ingredients"`
        Instructions	[]string	`json:"instructions"`
        PublishedAt		time.Time	`json:"publishedAt"`
    }
```

We first also define a **global variable** to store those recipes (for now we do not use a database):
```
    var recipes []Recipe
    func init() {						// Init function is called when the app is initialized
        recipes = make([]Recipe, 0) 	// Create the Recipe, with length of 0
    }
```

And we define how we will create a **new recipe**:
```
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
        recipes = append(recipes, recipe) // Append to the slice of recipes, the new recipe

        c.JSON(http.StatusOK, recipe)
    }
```

And for obtaining the recipes, we just have to pass to the body of request the **recipes** variable:
```
    // Define the Routes Handlers
    func GetSliceRecipes(c *gin.Context) {
        c.JSON(http.StatusOK, recipes)
    }
```
As in the recipes variable we stored all the **Recipes**.