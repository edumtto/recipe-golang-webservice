package html

import (
	"net/http"
	"text/template"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
)

// Renderer implements render.Interface to render HTML pages.
type Renderer struct{}

const (
	templatePath       = "../render/html/tmpl/"
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
func (Renderer) RenderRecipeList(w http.ResponseWriter, recipePreviews *[]domain.RecipePreview) {
	err := templates.ExecuteTemplate(w, listRecipeTemplate, recipePreviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipe renders a HTML page containing infomation about a specific recipe.
func (Renderer) RenderRecipe(w http.ResponseWriter, recipe *domain.Recipe) {
	err := templates.ExecuteTemplate(w, viewRecipeTemplate, recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipeEditor renders a HTML page containing a form to edit information from a specific recipe.
func (Renderer) RenderRecipeEditor(w http.ResponseWriter, recipeForm *domain.RecipeForm) {
	// TODO: Use recipeForm to render the available selectable options for category and difficulty
	err := templates.ExecuteTemplate(w, editRecipeTemplate, recipeForm.Recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderNewRecipeForm renders a HTML page containing an empty form to create a new recipe.
func (Renderer) RenderNewRecipeForm(w http.ResponseWriter, recipeForm *domain.RecipeForm) {
	err := templates.ExecuteTemplate(w, newRecipeTemplate, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
