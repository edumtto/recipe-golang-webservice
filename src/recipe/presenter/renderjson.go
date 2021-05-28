package presenter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Edu15/recipe-golang-webservice/src/recipe"
)

type selectableValues struct {
	Categories   []recipe.Category
	Difficulties []recipe.Difficulty
}

// ApiPresenter implements render.Interface to render JSON pages.
type ApiPresenter struct{}

// RenderRecipeList renders a JSON containing a list of recipes.
func (ApiPresenter) RenderRecipeList(w http.ResponseWriter, recipePreviews *[]recipe.Preview) {
	b, err := json.Marshal(recipePreviews)
	//err := json.ExecuteTemplate(w, listRecipeTemplate, recipePreviews)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipe renders a JSON containing infomation about a specific recipe.
func (ApiPresenter) RenderRecipe(w http.ResponseWriter, recipe *recipe.Entity) {
	b, err := json.Marshal(recipe)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipeEditor renders a JSON containing infomation about a specific recipe
// and the available values for selectable fields.
func (ApiPresenter) RenderRecipeEditor(w http.ResponseWriter, form *recipe.Form) {
	b, err := json.Marshal(form)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderNewRecipeForm renders a JSON containing the available values for selectable fields.
func (ApiPresenter) RenderNewRecipeForm(w http.ResponseWriter, form *recipe.Form) {
	selectableVals := selectableValues{
		Categories:   form.Categories,
		Difficulties: form.Difficulties,
	}

	b, err := json.Marshal(selectableVals)
	fmt.Fprint(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
