package recipe

import (
	"fmt"
	"net/http"
	"strconv"
)

type Service interface {
	List(w http.ResponseWriter, r *http.Request) (*[]Preview, error)
	View(w http.ResponseWriter, r *http.Request, id string) (*Entity, error)
	Edit(w http.ResponseWriter, r *http.Request, id string) (*Form, error)
	New(w http.ResponseWriter, r *http.Request) (*Form, error)
	Create(w http.ResponseWriter, r *http.Request) (int, error)
	Update(w http.ResponseWriter, r *http.Request, id string) error
	Delete(w http.ResponseWriter, r *http.Request, id string) error
}

// RecipeService is a http hander that provides use case methods to fetch and manipulate recipes from a repository.
type service struct {
	repo Repository
}

// NewService2 creates a new instance o RecipeService injecting a repository.
func NewService(repository Repository) Service {
	return &service{
		repo: repository,
	}
}

// List fetches a list of all recipes and present a formated result.
func (s service) List(w http.ResponseWriter, r *http.Request) (*[]Preview, error) {
	return s.repo.FetchPreviews()
}

// View fetches all information from recipe and present a formated result.
func (s service) View(w http.ResponseWriter, r *http.Request, id string) (*Entity, error) {
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
func (s service) Edit(w http.ResponseWriter, r *http.Request, id string) (*Form, error) {
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
func (s service) New(w http.ResponseWriter, r *http.Request) (*Form, error) {
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

func fetchFull(recipeID int, repo Repository) (*Entity, error) {
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

func fetchFormFieldValues(repo Repository) (*Form, error) {
	var recipeForm Form

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
