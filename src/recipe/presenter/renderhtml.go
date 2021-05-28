package presenter

import (
	"net/http"
	"text/template"

	"github.com/Edu15/recipe-golang-webservice/src/recipe"
)

// WebPresenter implements render.Interface to render HTML pages.
type WebPresenter struct{}

const (
	templatePath       = "../recipe/presenter/webtmpl/"
	listRecipeTemplate = "recipe-list.html"
	viewRecipeTemplate = "view-recipe.html"
	editRecipeTemplate = "edit-recipe.html"
	newRecipeTemplate  = "new-recipe.html"
)

var templates = template.Must(
	template.ParseFiles(
		templatePath+listRecipeTemplate,
		templatePath+viewRecipeTemplate,
		templatePath+editRecipeTemplate,
		templatePath+newRecipeTemplate,
	),
)

// RenderRecipeList renders a HTML page containing a list of recipes.
func (WebPresenter) RenderRecipeList(w http.ResponseWriter, recipePreviews *[]recipe.Preview) {
	err := templates.ExecuteTemplate(w, listRecipeTemplate, recipePreviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipe renders a HTML page containing infomation about a specific recipe.
func (WebPresenter) RenderRecipe(w http.ResponseWriter, recipe *recipe.Entity) {
	err := templates.ExecuteTemplate(w, viewRecipeTemplate, recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipeEditor renders a HTML page containing a form to edit information from a specific recipe.
func (WebPresenter) RenderRecipeEditor(w http.ResponseWriter, form *recipe.Form) {
	// TODO: Use recipeForm to render the available selectable options for category and difficulty
	err := templates.ExecuteTemplate(w, editRecipeTemplate, form.Recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderNewRecipeForm renders a HTML page containing an empty form to create a new recipe.
func (WebPresenter) RenderNewRecipeForm(w http.ResponseWriter, form *recipe.Form) {
	err := templates.ExecuteTemplate(w, newRecipeTemplate, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
