package service

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Edu15/recipe-golang-webservice/domain"
	"github.com/Edu15/recipe-golang-webservice/repository"
)

// RecipeService provides use case methods to fetch and manipulate recipes from a repository.
type RecipeService struct {
	repo repository.Interface
}

var templates = template.Must(template.ParseFiles("./tmpl/edit-recipe.html", "./tmpl/view-recipe.html", "./tmpl/recipe-list.html", "./tmpl/new-recipe.html"))

// func main() {
// 	fmt.Println("ok")
// }

// NewRecipeService creates a new instance o RecipeService injecting a repository.
func NewRecipeService() *RecipeService {
	repository := repository.NewRepository()
	return &RecipeService{
		repo: repository,
	}
}

// ListRecipes fetches a list of all recipes and present a formated result.
func (service *RecipeService) ListRecipes(w http.ResponseWriter, r *http.Request, id string) {
	recipePreviews, err := service.repo.FetchRecipePreviews(w, r)
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

func (service *RecipeService) ViewRecipe(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFullRecipe(recipeID, service.repo)
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

func (service *RecipeService) EditRecipe(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFullRecipe(recipeID, service.repo)
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

func (service *RecipeService) NewRecipe(w http.ResponseWriter, r *http.Request, id string) {
	renderTemplate(w, "new-recipe", nil)
}

func (service *RecipeService) CreateRecipe(w http.ResponseWriter, r *http.Request, id string) {
	insertedID, err := service.repo.InsertRecipe(w, r)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("/view/%d", insertedID)
	http.Redirect(w, r, url, http.StatusFound)
}

func (service *RecipeService) UpdateRecipe(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := service.repo.FetchRecipe(recipeID)
	if recipe == nil {
		panic(err)
	}

	err = service.repo.UpdateRecipe(w, r, recipeID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/view/"+id, http.StatusFound)
}

func (service *RecipeService) DeleteRecipe(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.repo.RemoveRecipe(w, r, recipeID)
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

func fetchFullRecipe(recipeID int, repo repository.Interface) (*domain.Recipe, error) {
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
