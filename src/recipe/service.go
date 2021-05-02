package recipe

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Edu15/recipe-golang-webservice/src/database"
	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/presenter/html"
	"github.com/Edu15/recipe-golang-webservice/src/presenter/json"
)

// RecipeService provides use case methods to fetch and manipulate recipes from a repository.
type Service struct {
	repo     domain.Repository
	renderer domain.Render
	format   domain.ResponseFormat
}

// NewService creates a new instance o RecipeService injecting a repository.
func NewService(format domain.ResponseFormat) *Service {
	repository := database.NewRepository(database.Connect())

	var renderer domain.Render
	if format == domain.JSON {
		renderer = json.Renderer{}
	} else {
		renderer = html.Renderer{}
	}

	return &Service{
		repo:     repository,
		renderer: renderer,
		format:   format,
	}
}

// List fetches a list of all recipes and present a formated result.
func (service *Service) List(w http.ResponseWriter, r *http.Request, id string) {
	recipePreviews, err := service.repo.FetchRecipePreviews(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	service.renderer.RenderRecipeList(w, recipePreviews)
}

// View fetches all information from recipe and present a formated result.
func (service *Service) View(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFull(recipeID, service.repo)
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

// Edit fetches recipe and present a form to edit the stored recipe information.
func (service *Service) Edit(w http.ResponseWriter, r *http.Request, id string) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipe, err := fetchFull(recipeID, service.repo)
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, r)
		return
	}

	if recipe == nil {
		fmt.Fprintf(w, "Receita não encontrada!")
		return
	}

	recipeForm, err := fetchFormFieldValues(service.repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipeForm.Recipe = *recipe
	service.renderer.RenderRecipeEditor(w, recipeForm)
}

// New renders a form to input information for a new recipe.
func (service *Service) New(w http.ResponseWriter, r *http.Request, id string) {
	recipeForm, err := fetchFormFieldValues(service.repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	service.renderer.RenderNewRecipeForm(w, recipeForm)
}

// Create persists a specified new recipe on the database.
func (service *Service) Create(w http.ResponseWriter, r *http.Request, id string) {
	insertedID, err := service.repo.InsertRecipe(w, r)
	if err != nil {
		fmt.Println(err)
	}

	if service.format == domain.HTML {
		url := fmt.Sprintf("/recipes/%d", insertedID)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// Update updates all information from a altered recipe on the database.
func (service *Service) Update(w http.ResponseWriter, r *http.Request, id string) {
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

	if service.format == domain.HTML {
		http.Redirect(w, r, "/recipes/"+id, http.StatusFound)
	}
}

// Delete removes all information from a specified recipe from the database.
func (service *Service) Delete(w http.ResponseWriter, r *http.Request, id string) {
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

	if service.format == domain.HTML {
		http.Redirect(w, r, "/recipes", http.StatusFound)
	}
}

func fetchFull(recipeID int, repo domain.Repository) (*domain.Recipe, error) {
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

	dificulty, err := repo.FetchDifficulty(recipe.Difficulty.ID)
	if err != nil {
		return nil, err
	}
	recipe.Difficulty.Name = dificulty.Name

	return recipe, nil
}

func fetchFormFieldValues(repo domain.Repository) (*domain.RecipeForm, error) {
	var recipeForm domain.RecipeForm

	categories, err := repo.FetchCategories()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	recipeForm.Categories = *categories

	difficulties, err := repo.FetchDifficulties()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	recipeForm.Difficulties = *difficulties

	return &recipeForm, nil
}
