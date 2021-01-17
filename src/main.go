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

var validPath = regexp.MustCompile("^/(view|edit|create|update|new|delete|all)/?([a-zA-Z0-9]+)?$")
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
	responseFormat := domain.HTML
	if len(os.Args) > 1 && os.Args[1] == "json" {
		responseFormat = domain.JSON
	}
	recipeService = service.NewRecipeService(responseFormat)

	http.HandleFunc("/all/", makeHandler(recipeService.ListRecipes))
	http.HandleFunc("/view/", makeHandler(recipeService.ViewRecipe))
	http.HandleFunc("/edit/", makeHandler(recipeService.EditRecipe))
	http.HandleFunc("/new/", makeHandler(recipeService.NewRecipe))
	http.HandleFunc("/create/", makeHandler(recipeService.CreateRecipe))
	http.HandleFunc("/update/", makeHandler(recipeService.UpdateRecipe))
	http.HandleFunc("/delete/", makeHandler(recipeService.DeleteRecipe))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
