package presenter

import (
	"net/http"
	"text/template"

	"github.com/Edu15/recipe-golang-webservice/src/domain"
	"github.com/Edu15/recipe-golang-webservice/src/recipe"
)

// WebPresenter implements render.Interface to render HTML pages.
type webPresenter struct {
	template domain.WebTemplate
}

var templates *template.Template

func NewWebPresenter(template domain.WebTemplate) recipe.Presenter {
	parseTemplates(template)

	return &webPresenter{
		template: template,
	}
}

func parseTemplates(t domain.WebTemplate) {
	templates = template.Must(
		template.ParseFiles(
			t.Path+t.ListFilename,
			t.Path+t.ViewFilename,
			t.Path+t.EditFilename,
			t.Path+t.NewFilename,
		),
	)
}

// RenderRecipeList renders a HTML page containing a list of recipes.
func (p webPresenter) RenderRecipeList(w http.ResponseWriter, recipePreviews *[]recipe.Preview) {
	err := templates.ExecuteTemplate(w, p.template.ListFilename, recipePreviews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipe renders a HTML page containing infomation about a specific recipe.
func (p webPresenter) RenderRecipe(w http.ResponseWriter, recipe *recipe.Entity) {
	err := templates.ExecuteTemplate(w, p.template.ViewFilename, recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderRecipeEditor renders a HTML page containing a form to edit information from a specific recipe.
func (p webPresenter) RenderRecipeEditor(w http.ResponseWriter, form *recipe.Form) {
	// TODO: Use recipeForm to render the available selectable options for category and difficulty
	err := templates.ExecuteTemplate(w, p.template.EditFilename, form.Recipe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderNewRecipeForm renders a HTML page containing an empty form to create a new recipe.
func (p webPresenter) RenderNewRecipeForm(w http.ResponseWriter, form *recipe.Form) {
	err := templates.ExecuteTemplate(w, p.template.NewFilename, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
