package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/Edu15/recipe-golang-webservice/service"
)

var validPath = regexp.MustCompile("^/(view|edit|create|update|new|delete|home)/?([a-zA-Z0-9]+)?$")
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

func main() {
	recipeService = service.NewRecipeService()

	http.HandleFunc("/home/", makeHandler(recipeService.ListRecipes))
	http.HandleFunc("/view/", makeHandler(recipeService.ViewRecipe))
	http.HandleFunc("/edit/", makeHandler(recipeService.EditRecipe))
	http.HandleFunc("/new/", makeHandler(recipeService.NewRecipe))
	http.HandleFunc("/create/", makeHandler(recipeService.CreateRecipe))
	http.HandleFunc("/update/", makeHandler(recipeService.UpdateRecipe))
	http.HandleFunc("/delete/", makeHandler(recipeService.DeleteRecipe))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
