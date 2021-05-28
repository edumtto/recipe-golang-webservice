package recipe

import (
	"net/http"
)

// Interface interface
type Repository interface {
	FetchRecipe(recipeID int) (*Entity, error)
	FetchAuthor(ID int) (*Author, error)
	FetchCategory(ID int) (*Category, error)
	FetchDifficulty(ID int) (*Difficulty, error)
	FetchPreviews() (*[]Preview, error)
	UpdateRecipe(w http.ResponseWriter, r *http.Request, id int) error
	InsertRecipe(w http.ResponseWriter, r *http.Request) (int, error)
	RemoveRecipe(w http.ResponseWriter, r *http.Request, id int) error
	FetchCategories() (*[]Category, error)
	FetchDifficulties() (*[]Difficulty, error)
}

// Interface for rendering methods
type Presenter interface {
	RenderRecipeList(w http.ResponseWriter, previews *[]Preview)
	RenderRecipe(w http.ResponseWriter, recipe *Entity)
	RenderRecipeEditor(w http.ResponseWriter, form *Form)
	RenderNewRecipeForm(w http.ResponseWriter, Form *Form)
}
