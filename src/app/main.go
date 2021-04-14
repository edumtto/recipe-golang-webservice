package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/recipe"
)

var validPath = regexp.MustCompile("^/(recipes|edit|new|create|update|delete)/?([a-zA-Z0-9]+)?$")
var recipeService *recipe.Service

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
			recipeService.New(w, r, "")
		} else if recipeID != "" {
			recipeService.View(w, r, recipeID)
		} else {
			recipeService.List(w, r, "")
		}

	case http.MethodPost:
		recipeService.Create(w, r, "")

	case http.MethodPut:
		recipeService.Update(w, r, recipeID)

	case http.MethodDelete:
		recipeService.Delete(w, r, recipeID)

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
	responseFormat := domain.JSON
	if len(os.Args) > 1 && os.Args[1] == "html" {
		responseFormat = domain.HTML
	}
	recipeService = recipe.NewService(responseFormat)

	http.HandleFunc("/recipes/", recipesHandler)
	http.HandleFunc("/new/", makeHandler(recipeService.New))
	http.HandleFunc("/create/", makeHandler(recipeService.Create))
	http.HandleFunc("/edit/", makeHandler(recipeService.Edit))
	http.HandleFunc("/update/", makeHandler(recipeService.Update))
	http.HandleFunc("/delete/", makeHandler(recipeService.Delete))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
