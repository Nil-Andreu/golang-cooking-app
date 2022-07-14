# Cooking App
Will create a cooking application with Golang.

In this application, you will be able to handle recipes.

## Defining Recipe Structure
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

## Define POST request
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
        recipe.ID = xid.New().String()   // Create a unique identifier, like: "cb7s7ne49b3mr2q8ci00"
        recipe.PublishedAt = time.Now()
        recipes = append(recipes, recipe) // Append to the slice of recipes, the new recipe

        c.JSON(http.StatusOK, recipe)
    }
```

## Define GET request
And for obtaining the recipes, we just have to pass to the body of request the **recipes** variable:
```
    // Define the Routes Handlers
    func 
    (c *gin.Context) {
        c.JSON(http.StatusOK, recipes)
    }
```
As in the recipes variable we stored all the **Recipes**.

## Data Persistency with JSON file
We will lose the list of the recipes when we re-initialize the project. To have data persistency, we could store the recipes in a .json file and then when initializing the application to read this file.
```
    // Function init is ran when the app is initialized
    func init() { // Init function is called when the app is initialized
        recipes = make([]Recipe, 0) // Create the Recipe, with length of 0
        file, _ := ioutil.ReadFile("recipes.json") // Read the file
        _ = json.Unmarshal([]byte(file), &recipes)
    }
```

And not only reading, but we might want to also write to this file each time we make a POST request.
```
    ...
    // We write the new information
	file, _ := json.MarshalIndent(recipes, "", "")   // Write recipes as json structure in bytes
	_ = ioutil.WriteFile("recipes.json", file, 0644) // store it
```
So in each post request, we are storing the recipes data structure in the recipes.json file (which is read when we initialize the app).

## Updating Data