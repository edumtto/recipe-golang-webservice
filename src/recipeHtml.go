package main

import (
	"fmt"
	"log"
	"net/http"
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

/*
GET /new to fech a form to create a new recipe
GET /create to fech a form to create a new recipe
GET /edit to fech a form to create a new recipe
GET /update/{recipeId} to update an existing recipe
GET /delete/{recipeId} to delete a recipe
*/

func main() {
	recipeService = recipe.NewService(domain.HTML)

	http.HandleFunc("/new/", makeHandler(recipeService.New))
	http.HandleFunc("/create/", makeHandler(recipeService.Create))
	http.HandleFunc("/edit/", makeHandler(recipeService.Edit))
	http.HandleFunc("/update/", makeHandler(recipeService.Update))
	http.HandleFunc("/delete/", makeHandler(recipeService.Delete))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
