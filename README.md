# Cooking App
Will create a cooking application with Golang.

In this application, you will be able to handle recipes.

## Defining Recipe Structure
The object of the recipes is the following:
```
    type Recipe struct {
        ID string `json:"id"`	
        Name string `json:"name"`
        Tags []string `json:"ingredients"`
        Instructions []string `json:"instructions"`
        PublishedAt	time.Time `json:"publishedAt"`
    }
```

We first also define a **global variable** to store those recipes (for now we do not use a database):
```
    var recipes []Recipe
    func init() {  // Init function is called when the app is initialized
        recipes = make([]Recipe, 0)  // Create the Recipe, with length of 0
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
        recipe.ID = xid.New().String()  // Create a unique identifier, like: "cb7s7ne49b3mr2q8ci00"
        recipe.PublishedAt = time.Now()
        recipes = append(recipes, recipe)  // Append to the slice of recipes, the new recipe

        c.JSON(http.StatusOK, recipe)
    }
```

## Define GET request
And for obtaining the recipes, we just have to pass to the body of request the **recipes** variable:
```
    // Define the Routes Handlers
    func GetSliceRecipes(c *gin.Context) {
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
	file, _ := json.MarshalIndent(recipes, "", "")  // Write recipes as json structure in bytes
	_ = ioutil.WriteFile("recipes.json", file, 0644) // store it
```
So in each post request, we are storing the recipes data structure in the recipes.json file (which is read when we initialize the app).

## Updating Data
For updating the recipe, we use the PUT method and we pass in the url the id of the recipe we want to update.
Then, if this id is on one of the recipes, we replace that with the body of the request.
```
    // To update the Recipe
    func UpdateRecipeHandler(c *gin.Context) {
        // Obtain first which is the parameter
        id := c.Param("id")

        // Initialize the recipe
        var recipe Recipe

        // Bind the recipe of theb ody of the request to json
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
            response := fmt.Sprintf("Index %s not found", id)
            c.JSON(http.StatusBadRequest, gin.H{
                "error": response,
            })
        }

        // If we found it, replace the recipe at that position
        recipes[index] = recipe

        // And we now need to re store it on the json file
        file, _ := json.MarshalIndent(recipes, "", "") // Write recipes as json structure in bytes
        _ = ioutil.WriteFile("recipes.json", file, 0644)

        c.JSON(http.StatusOK, recipe)
    }
```

## Deletetion of the Recipe
For the deletion, we follow the same process as the update, where we find on which positon on the recipes slice it is the one we want. And then we delete it with:
```
    recipes = append(recipes[:index], recipes[index+1:]...)
```
Where this is the Golang way to remove an item at a certain position.

## Search for Tags
Another interesting thing we might want to do is to search for recipes that have an specific tag.
This will be a GET method, where the tag will be in the url but as *?tag=* query parameter.
```
    func SearchRecipeHandler(c *gin.Context) {
        // Obtain first the query parameter
	    Tag := c.Query("tag")
	    resultRecipes := make([]Recipe, 0)

	    for recipeIndex := 0; recipeIndex < len(recipes); recipeIndex ++ {
            // And for this recipe, we will look for all the tags that it has
            recipe := recipes[recipeIndex]

            for tagIndex := 0; tagIndex < len(recipe.Tag); tagIndex ++ {
                // In the case that those are equal
                if strings.EqualFold(recipe.Tag[tagIndex], Tag) {
                    resultRecipes = append(resultRecipes, recipe)
                }
            }

            // In the case there is no match, handle this situation
            if len(resultRecipes) == 0 {
                c.JSON(http.StatusOK, gin.H {
                    "message" : "Not found",
                })
            } else {
                c.JSON(http.StatusOK, resultRecipes)
            }
        }
    }
```
For comparing the two strings, we use *strings.EqualFold*, so we do not worry about if one letter is in upper or lower case.

## Documentation
Once we have defined the APIs, we could create a documentation.
The OpenAPI Specifications is an API description format which tells:
- General information about the API
- Available paths and HTTP methods for each path
- Expected inputs and responses for each path

For this, we can download *go-swagger*:
```
    brew tap go-swagger/go-swagger
    brew install go-swagger
```

The Swagger Metadata we would have to provide:
- **Schemes**: transfer prototols supported in the API (http and https)
- **Host**: the host where the api is served (localhost:6000)
- **BasePath**: default base path for the api (/v1) for example
- **Version**: current version of the API 
- **Contact**: owner/author of the API
- **Consumes**: list of MIME type values, one per line for the content the API receives (f.e. application/json)
- **Produces**: list of default MIME types the application sends, one per line

So for this, we have to add this infroamtion at the top of the *main.go* file.
```
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
```

And once we have defined all of this, we run the following command: *swagger generate spec -o ./swagger.json*.
This will generate the *swagger.json* file with all the information of our API.
```
    {
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "schemes": [
        "http"
    ],
    "swagger": "2.0",
    "info": {
        "description": "Recipes\nA simple API for adding, removing, reviewing and updating recipes.",
        "contact": {
        "name": "Nil Andreu",
        "email": "nilandreug@gmail.com"
        },
        "version": "1.0.0"
    },
    "host": "localhost:6000",
    "basePath": "/",
    "paths": {}
    }
```
And to serve it, we would put on the command line: *swagger serve ./swagger.json*.
Which will show you a simple UI for the documentation of the API.

In the case we wanted the UI of the Swagger: *swagger serve -F swagger ./swagger*.

But now we do not have information about the different endpoints. This has to be added by us. 
For each one, we need to add:
- **Summary**: short description of what it does
- **Responses**: list of possible responses that are possible when executing this operation
- **Parameters**: parameters expected for the operation
- **Consumes**: list of default MIME type values, one per line
- **Produces**: list of default MIME type values, one per line

So for example in the API for obtaining all the recipes:
```
    // swagger:operation GET /recipes recipes listRecipes
    // Returns list of recipes
    // ---
    // produces:
    // - application/json
    // responses:
    //  '200':
    //	  description: Successful operation
```
We do not make a tab, but use spaces.

And we could do also for the other endpoints that we have:
```
    // swagger:operation POST /recipes recipe
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
```
Where note that for example in the tag parameter:
```
    // - name: tag
    //   in: body
    //   description: Tags for the recipe
    //   require: true
    //   schema: 
    //    type: array
    //    items:
    //     type: string
    //    example: ["tag1", "tag2"]
```
We define the following:
- the name
- where is it (body, path, ...)
- description
- if it is required or not
- schema: the type and the items inside of it