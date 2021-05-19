package domain

import (
	"net/http"
)

// Interface interface
type Repository interface {
	FetchRecipe(recipeID int) (*Recipe, error)
	FetchAuthor(ID int) (*RecipeAuthor, error)
	FetchCategory(ID int) (*RecipeCategory, error)
	FetchDifficulty(ID int) (*RecipeDifficulty, error)
	FetchRecipePreviews() (*[]RecipePreview, error)
	UpdateRecipe(w http.ResponseWriter, r *http.Request, id int) error
	InsertRecipe(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveRecipe(w http.ResponseWriter, r *http.Request, id int) error
	FetchCategories() (*[]RecipeCategory, error)
	FetchDifficulties() (*[]RecipeDifficulty, error)
}

// Interface for rendering methods
type Render interface {
	RenderRecipeList(w http.ResponseWriter, recipePreviews *[]RecipePreview)
	RenderRecipe(w http.ResponseWriter, recipe *Recipe)
	RenderRecipeEditor(w http.ResponseWriter, recipeForm *RecipeForm)
	RenderNewRecipeForm(w http.ResponseWriter, recipeForm *RecipeForm)
}
