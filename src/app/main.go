package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Edu15/recipe-golang-webservice/src/database"
	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/recipe"
	"github.com/Edu15/recipe-golang-webservice/src/recipe/presenter"
)

var validWebPath = regexp.MustCompile("^/recipes/(all|edit|new|create|update|delete)?/?([a-zA-Z0-9]+)?$")
var validApiPath = regexp.MustCompile("^/api/recipes/?([a-zA-Z0-9]+)?$")
var webController recipe.Controller
var apiController recipe.Controller

var databaseConf = database.DatabaseConf{
	Host:     "localhost",
	Port:     5432,
	User:     "edu",
	Password: "1234",
	DbName:   "recipes_db",
}

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
			apiController.View(w, r, recipeID)
		} else {
			apiController.List(w, r)
		}

	case http.MethodPost:
		apiController.Create(w, r)

	case http.MethodPut:
		apiController.Update(w, r, recipeID)

	case http.MethodDelete:
		apiController.Delete(w, r, recipeID)

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
	repository := recipe.NewRepository(database.Connect(databaseConf))
	service := recipe.NewService(repository)
	webController = recipe.NewController(service, presenter.WebPresenter{}, domain.HTML)
	apiController = recipe.NewController(service, presenter.ApiPresenter{}, domain.JSON)

	http.HandleFunc("/api/recipes/", recipesApiHandler)
	http.HandleFunc("/api/recipes/form/", apiController.New)

	http.HandleFunc("/recipes/", makeWebHandler(webController.View))
	http.HandleFunc("/recipes/all/", webController.List)
	http.HandleFunc("/recipes/new/", webController.New)
	http.HandleFunc("/recipes/create/", webController.Create)
	http.HandleFunc("/recipes/edit/", makeWebHandler(webController.Edit))
	http.HandleFunc("/recipes/update/", makeWebHandler(webController.Update))
	http.HandleFunc("/recipes/delete/", makeWebHandler(webController.Delete))

	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
