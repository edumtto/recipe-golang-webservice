package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"text/template"

	"github.com/Edu15/recipe-golang-webservice/domain"
	"github.com/Edu15/recipe-golang-webservice/repository"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "edu"
	password = "1234"
	dbname   = "recipes_db"
)

var templates = template.Must(template.ParseFiles("tmpl/edit-recipe.html", "tmpl/view-recipe.html", "tmpl/recipe-list.html", "tmpl/new-recipe.html"))
var validPath = regexp.MustCompile("^/(view|edit|create|update|new|delete|home)/?([a-zA-Z0-9]+)?$")
var repo repository.Interface

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

func listRecipesHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipePreviews, err := repo.FetchRecipePreviews(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = templates.ExecuteTemplate(w, "recipe-list.html", recipePreviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFullRecipe(recipeID)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	if recipe == nil {
		fmt.Fprintf(w, "Receita não encontrada!")
		return
	}

	renderTemplate(w, "view-recipe", recipe)
}

func editRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFullRecipe(recipeID)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	if recipe == nil {
		fmt.Fprintf(w, "Receita não encontrada!")
		return
	}

	renderTemplate(w, "edit-recipe", recipe)
}

func newRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	renderTemplate(w, "new-recipe", nil)
}

func createRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	insertedID, err := repo.InsertRecipe(w, r)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("/view/%d", insertedID)
	http.Redirect(w, r, url, http.StatusFound)
}

func updateRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := repo.FetchRecipe(recipeID)
	if recipe == nil {
		panic(err)
	}

	err = repo.UpdateRecipe(w, r, recipeID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/view/"+id, http.StatusFound)
}

func deleteRecipeHandler(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = repo.RemoveRecipe(w, r, recipeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, recipe *domain.Recipe) {
	err := templates.ExecuteTemplate(w, tmpl+".html", recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchFullRecipe(recipeID int) (*domain.Recipe, error) {
	recipe, err := repo.FetchRecipe(recipeID)
	if err != nil {
		return nil, err
	}

	author, err := repo.FetchAuthor(recipe.Author.ID)
	if err != nil {
		return nil, err
	}
	recipe.Author.Name = author.Name

	category, err := repo.FetchCategory(recipe.Category.ID)
	if err != nil {
		return nil, err
	}
	recipe.Category.Name = category.Name

	dificulty, err := repo.FetchDificulty(recipe.Dificulty.ID)
	if err != nil {
		return nil, err
	}
	recipe.Dificulty.Name = dificulty.Name

	return recipe, nil
}

func main() {
	repo = repository.NewRepository()

	http.HandleFunc("/home/", makeHandler(listRecipesHandler))
	http.HandleFunc("/view/", makeHandler(viewRecipeHandler))
	http.HandleFunc("/edit/", makeHandler(editRecipeHandler))
	http.HandleFunc("/new/", makeHandler(newRecipeHandler))
	http.HandleFunc("/create/", makeHandler(createRecipeHandler))
	http.HandleFunc("/update/", makeHandler(updateRecipeHandler))
	http.HandleFunc("/delete/", makeHandler(deleteRecipeHandler))
	fmt.Println("Service is running.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
