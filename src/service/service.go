package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/render"
	"github.com/Edu15/recipe-golang-webservice/src/repository"
)

// RecipeService provides use case methods to fetch and manipulate recipes from a repository.
type RecipeService struct {
	repo     repository.Interface
	renderer render.Interface
}

// NewRecipeService creates a new instance o RecipeService injecting a repository.
func NewRecipeService() *RecipeService {
	repository := repository.NewRepository()
	return &RecipeService{
		repo:     repository,
		renderer: render.HTMLRenderer{},
	}
}

// ListRecipes fetches a list of all recipes and present a formated result.
func (service *RecipeService) ListRecipes(w http.ResponseWriter, r *http.Request, id string) {
	recipePreviews, err := service.repo.FetchRecipePreviews(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	service.renderer.RenderRecipeList(w, recipePreviews)
}

// ViewRecipe fetches all information from recipe and present a formated result.
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

	service.renderer.RenderRecipe(w, recipe)
}

// EditRecipe fetches recipe and present a form to edit the stored recipe information.
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

	service.renderer.RenderRecipeEditor(w, recipe)
}

// NewRecipe renders a form to input information for a new recipe.
func (service *RecipeService) NewRecipe(w http.ResponseWriter, r *http.Request, id string) {
	service.renderer.RenderNewRecipeForm(w)
}

// CreateRecipe persists a specified new recipe on the database.
func (service *RecipeService) CreateRecipe(w http.ResponseWriter, r *http.Request, id string) {
	insertedID, err := service.repo.InsertRecipe(w, r)
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("/view/%d", insertedID)
	http.Redirect(w, r, url, http.StatusFound)
}

// UpdateRecipe updates all information from a altered recipe on the database.
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

// DeleteRecipe removes all information from a specified recipe from the database.
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
