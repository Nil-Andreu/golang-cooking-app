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

We first also define a global variable to store those recipes (for now we do not use a database):
```
    var recipes []Recipe
    func init() {						// Init function is called when the app is initialized
        recipes = make([]Recipe, 0) 	// Create the Recipe, with length of 0
    }
```
