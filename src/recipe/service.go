package recipe

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
)

type Service interface {
	List(w http.ResponseWriter, r *http.Request) (*[]domain.RecipePreview, error)
	View(w http.ResponseWriter, r *http.Request, id string) (*domain.Recipe, error)
	Edit(w http.ResponseWriter, r *http.Request, id string) (*domain.RecipeForm, error)
	New(w http.ResponseWriter, r *http.Request) (*domain.RecipeForm, error)
	Create(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request, id string) error
	Delete(w http.ResponseWriter, r *http.Request, id string) error
}

// RecipeService is a http hander that provides use case methods to fetch and manipulate recipes from a repository.
type service struct {
	repo     domain.Repository
	renderer domain.Render
	format   domain.ResponseFormat
}

// NewService2 creates a new instance o RecipeService injecting a repository.
func NewService(repository domain.Repository, renderer domain.Render, format domain.ResponseFormat) Service {
	return &service{
		repo:     repository,
		renderer: renderer,
		format:   format,
	}
}

// List fetches a list of all recipes and present a formated result.
func (s service) List(w http.ResponseWriter, r *http.Request) (*[]domain.RecipePreview, error) {
	return s.repo.FetchRecipePreviews()
}

// View fetches all information from recipe and present a formated result.
func (s service) View(w http.ResponseWriter, r *http.Request, id string) (*domain.Recipe, error) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	recipe, err := fetchFull(recipeID, s.repo)
	if err != nil {
		return nil, err
	}

	if recipe == nil {
		return nil, fmt.Errorf("Receita com id %snão encontrada", id)
	}

	return recipe, nil
}

// Edit fetches recipe and present a form to edit the stored recipe information.
func (s service) Edit(w http.ResponseWriter, r *http.Request, id string) (*domain.RecipeForm, error) {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	recipe, err := fetchFull(recipeID, s.repo)
	if err != nil {
		return nil, err
	}

	if recipe == nil {
		return nil, fmt.Errorf("Receita com id %snão encontrada", id)
	}

	recipeForm, err := fetchFormFieldValues(s.repo)
	if err != nil {
		return nil, err
	}

	recipeForm.Recipe = *recipe
	return recipeForm, nil
}

// New renders a form to input information for a new recipe.
func (s service) New(w http.ResponseWriter, r *http.Request) (*domain.RecipeForm, error) {
	return fetchFormFieldValues(s.repo)
}

// Create persists a specified new recipe on the database.
func (s service) Create(w http.ResponseWriter, r *http.Request) (int, error) {
	return s.repo.InsertRecipe(w, r)
}

// Update updates all information from a altered recipe on the database.
func (s service) Update(w http.ResponseWriter, r *http.Request, id string) error {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	recipe, err := s.repo.FetchRecipe(recipeID)
	if recipe == nil {
		return err
	}

	return s.repo.UpdateRecipe(w, r, recipeID)
}

// Delete removes all information from a specified recipe from the database.
func (s service) Delete(w http.ResponseWriter, r *http.Request, id string) error {
	recipeID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return s.repo.RemoveRecipe(w, r, recipeID)
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
