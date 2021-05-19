package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Edu15/recipe-golang-webservice/src/database"
	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/presenter/html"
	"github.com/Edu15/recipe-golang-webservice/src/presenter/json"
	"github.com/Edu15/recipe-golang-webservice/src/recipe"
)

var validWebPath = regexp.MustCompile("^/recipes/(all|edit|new|create|update|delete)?/?([a-zA-Z0-9]+)?$")
var validApiPath = regexp.MustCompile("^/api/recipes/?([a-zA-Z0-9]+)?$")
var webService *recipe.Service
var apiService *recipe.Service

func makeWebHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validWebPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fmt.Println("Invalid web path!")
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func recipesApiHandler(w http.ResponseWriter, r *http.Request) {
	m := validApiPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		fmt.Println("Invalid api path!")
		http.NotFound(w, r)
		return
	}

	var recipeID = ""
	if len(m) > 1 {
		recipeID = m[1]
	}

	switch r.Method {
	case http.MethodGet:
		if recipeID != "" {
			apiService.View(w, r, recipeID)
		} else {
			apiService.List(w, r)
		}

	case http.MethodPost:
		apiService.Create(w, r)

	case http.MethodPut:
		apiService.Update(w, r, recipeID)

	case http.MethodDelete:
		apiService.Delete(w, r, recipeID)

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
	repository := database.NewRepository(database.Connect())
	webService = recipe.NewService(repository, html.Renderer{}, domain.HTML)
	apiService = recipe.NewService(repository, json.Renderer{}, domain.JSON)

	http.HandleFunc("/api/recipes/", recipesApiHandler)
	http.HandleFunc("/api/recipes/form/", apiService.New)

	http.HandleFunc("/recipes/", makeWebHandler(webService.View))
	http.HandleFunc("/recipes/all/", webService.List)
	http.HandleFunc("/recipes/new/", webService.New)
	http.HandleFunc("/recipes/create/", webService.Create)
	http.HandleFunc("/recipes/edit/", makeWebHandler(webService.Edit))
	http.HandleFunc("/recipes/update/", makeWebHandler(webService.Update))
	http.HandleFunc("/recipes/delete/", makeWebHandler(webService.Delete))

	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
