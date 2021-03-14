package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/service"
)

var validPath = regexp.MustCompile("^/(recipes|edit|new|create|update|delete)/?([a-zA-Z0-9]+)?$")
var recipeService *service.RecipeService

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fmt.Println("Invalid path!")
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func recipesHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		fmt.Println("Invalid path!")
		http.NotFound(w, r)
		return
	}

	recipeID := m[2]

	switch r.Method {
	case http.MethodGet:
		if recipeID == "form" {
			recipeService.NewRecipe(w, r, "")
		} else if recipeID != "" {
			recipeService.ViewRecipe(w, r, recipeID)
		} else {
			recipeService.ListRecipes(w, r, "")
		}

	case http.MethodPost:
		recipeService.CreateRecipe(w, r, "")

	case http.MethodPut:
		recipeService.UpdateRecipe(w, r, recipeID)

	case http.MethodDelete:
		recipeService.DeleteRecipe(w, r, recipeID)

	default:
		http.NotFound(w, r)
	}
}

/*

POST /recipes to add a new recipe
GET /recipes to fetch all existing recipes
GET /recipes/form to fech a form to create a new recipe
GET /recipes/{recipeId} to fetch a single recipe using its ID
PUT /recipes/{recipeId} to update an existing recipe
DELETE /recipes/{itemId} to delete a recipe

// Only for HTML forms:
GET /new to fech a form to create a new recipe
GET /create to fech a form to create a new recipe
GET /edit to fech a form to create a new recipe
GET /update/{recipeId} to update an existing recipe
GET /delete/{recipeId} to delete a recipe

*/

func main() {
	responseFormat := domain.HTML
	if len(os.Args) > 1 && os.Args[1] == "json" {
		responseFormat = domain.JSON
	}
	recipeService = service.NewRecipeService(responseFormat)

	http.HandleFunc("/recipes/", recipesHandler)
	http.HandleFunc("/new/", makeHandler(recipeService.NewRecipe))
	http.HandleFunc("/create/", makeHandler(recipeService.CreateRecipe))
	http.HandleFunc("/edit/", makeHandler(recipeService.EditRecipe))
	http.HandleFunc("/update/", makeHandler(recipeService.UpdateRecipe))
	http.HandleFunc("/delete/", makeHandler(recipeService.DeleteRecipe))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
