package main

import (
	"fmt"
	"log"
	"net/http"
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
*/

func main() {
	recipeService = service.NewRecipeService(domain.JSON)

	http.HandleFunc("/recipes/", recipesHandler)
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
