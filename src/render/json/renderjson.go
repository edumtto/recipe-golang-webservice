package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
)

type selectableValues struct {
	Categories   []domain.RecipeCategory
	Difficulties []domain.RecipeDifficulty
}

// Renderer implements render.Interface to render JSON pages.
type Renderer struct{}

// RenderRecipeList renders a JSON containing a list of recipes.
func (Renderer) RenderRecipeList(w http.ResponseWriter, recipePreviews *[]domain.RecipePreview) {
	b, err := json.Marshal(recipePreviews)
	//err := json.ExecuteTemplate(w, listRecipeTemplate, recipePreviews)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipe renders a JSON containing infomation about a specific recipe.
func (Renderer) RenderRecipe(w http.ResponseWriter, recipe *domain.Recipe) {
	b, err := json.Marshal(recipe)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipeEditor renders a JSON containing infomation about a specific recipe
// and the available values for selectable fields.
func (Renderer) RenderRecipeEditor(w http.ResponseWriter, recipeForm *domain.RecipeForm) {
	b, err := json.Marshal(recipeForm)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderNewRecipeForm renders a JSON containing the available values for selectable fields.
func (Renderer) RenderNewRecipeForm(w http.ResponseWriter, recipeForm *domain.RecipeForm) {
	selectableVals := selectableValues{
		Categories:   recipeForm.Categories,
		Difficulties: recipeForm.Difficulties,
	}

	b, err := json.Marshal(selectableVals)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
